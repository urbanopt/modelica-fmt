// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestModelicaEscapeRoundTrip(t *testing.T) {
	a := require.New(t)
	cases := []string{
		`<a href="modelica://Foo.Bar">Foo</a>`,
		"line with a backslash \\ and a quote \"",
		"plain text no escapes",
		"tab\tand\nnewline",
	}
	for _, decoded := range cases {
		escaped := escapeModelicaString(decoded)
		a.Equal(decoded, unescapeModelicaString(escaped), "round trip should be lossless")
	}
}

func TestUnescapeModelicaString(t *testing.T) {
	a := require.New(t)
	a.Equal(`href="x"`, unescapeModelicaString(`href=\"x\"`))
	a.Equal(`a\b`, unescapeModelicaString(`a\\b`))
	a.Equal("a\nb", unescapeModelicaString(`a\nb`))
}

func TestIsHTMLContent(t *testing.T) {
	a := require.New(t)
	a.True(isHTMLContent("<html><p>hi</p></html>"))
	a.True(isHTMLContent("\n  <html>\nhi\n</html>"))
	a.True(isHTMLContent("<HTML>hi</HTML>"))
	a.True(isHTMLContent(`<html lang="en">hi</html>`))
	a.False(isHTMLContent("just some text"))
	a.False(isHTMLContent("<p>not wrapped in html</p>"))
	a.False(isHTMLContent(""))
}

func TestFormatHTMLDocStringBasic(t *testing.T) {
	a := require.New(t)
	out, ok := formatHTMLDocString("<html><p>hi</p></html>", 0)
	a.True(ok)
	a.Equal("<html>\n  <p>\n    hi\n  </p>\n</html>", out)
}

func TestFormatHTMLDocStringBaseIndent(t *testing.T) {
	a := require.New(t)
	out, ok := formatHTMLDocString("<html>hi</html>", 2)
	a.True(ok)
	a.Equal("    <html>\n      hi\n    </html>", out)
}

func TestFormatHTMLDocStringIdempotent(t *testing.T) {
	a := require.New(t)
	src := "<html>\n<p>Hello &amp; world <b>bold</b> and <a href=\"x\">link</a>.</p>\n<br/>\n</html>"
	once, ok := formatHTMLDocString(src, 1)
	a.True(ok)
	twice, ok := formatHTMLDocString(once, 1)
	a.True(ok)
	a.Equal(once, twice, "formatting should be idempotent")
}

func TestFormatHTMLDocStringPreservesEntities(t *testing.T) {
	a := require.New(t)
	out, ok := formatHTMLDocString("<html><p>a &amp; b &lt; c &nbsp; d</p></html>", 0)
	a.True(ok)
	a.Contains(out, "a &amp; b &lt; c &nbsp; d")
}

func TestFormatHTMLDocStringPreservesAttributesWithGt(t *testing.T) {
	a := require.New(t)
	// A '>' inside a quoted attribute value must not split the tag.
	out, ok := formatHTMLDocString(`<html><span title="a > b">x</span></html>`, 0)
	a.True(ok)
	a.Contains(out, `<span title="a > b">`)
}

func TestFormatHTMLDocStringPreservesPre(t *testing.T) {
	a := require.New(t)
	src := "<html>\n<pre>\n  keep   this\n    exactly\n</pre>\n</html>"
	out, ok := formatHTMLDocString(src, 0)
	a.True(ok)
	// The whitespace-sensitive <pre> content must be preserved verbatim.
	a.Contains(out, "<pre>\n  keep   this\n    exactly\n</pre>")
}

func TestFormatHTMLDocStringMalformedBails(t *testing.T) {
	a := require.New(t)
	// Unbalanced tags should cause a graceful bail (ok == false).
	_, ok := formatHTMLDocString("<html><p>unclosed paragraph</html>", 0)
	a.False(ok)

	_, ok = formatHTMLDocString("<html><b>mismatched</i></html>", 0)
	a.False(ok)

	_, ok = formatHTMLDocString("<html>never closed", 0)
	a.False(ok)
}

func TestMaybeFormatHTMLString(t *testing.T) {
	a := require.New(t)

	// Non-HTML string: returned unchanged.
	orig := `"just plain text"`
	out, ok := maybeFormatHTMLString(orig, 0)
	a.False(ok)
	a.Equal(orig, out)

	// HTML string with escaped quotes: formatted and re-escaped losslessly.
	in := `"<html><a href=\"x\">y</a></html>"`
	out, ok = maybeFormatHTMLString(in, 0)
	a.True(ok)
	a.Equal("\"\n<html>\n  <a href=\\\"x\\\">\n    y\n  </a>\n</html>\"", out)

	// Malformed HTML string: returned unchanged.
	bad := `"<html><p>no close"`
	out, ok = maybeFormatHTMLString(bad, 0)
	a.False(ok)
	a.Equal(bad, out)
}
