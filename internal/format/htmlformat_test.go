// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package format

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
	out, err := formatHTMLDocString("<html><p>hi</p></html>", 0)
	a.NoError(err)
	a.Equal("<html>\n  <p>hi</p>\n</html>", out)
}

func TestFormatHTMLDocStringBaseIndent(t *testing.T) {
	a := require.New(t)
	out, err := formatHTMLDocString("<html>hi</html>", 2)
	a.NoError(err)
	a.Equal("    <html>\n      hi\n    </html>", out)
}

func TestFormatHTMLDocStringIdempotent(t *testing.T) {
	a := require.New(t)
	src := "<html>\n<p>Hello &amp; world <b>bold</b> and <a href=\"x\">link</a>.</p>\n<br/>\n</html>"
	once, err := formatHTMLDocString(src, 1)
	a.NoError(err)
	twice, err := formatHTMLDocString(once, 1)
	a.NoError(err)
	a.Equal(once, twice, "formatting should be idempotent")
}

func TestFormatHTMLDocStringPreservesEntities(t *testing.T) {
	a := require.New(t)
	out, err := formatHTMLDocString("<html><p>a &amp; b &lt; c &nbsp; d</p></html>", 0)
	a.NoError(err)
	a.Contains(out, "a &amp; b &lt; c &nbsp; d")
}

func TestFormatHTMLDocStringPreservesAttributesWithGt(t *testing.T) {
	a := require.New(t)
	// A '>' inside a quoted attribute value must not split the tag.
	out, err := formatHTMLDocString(`<html><span title="a > b">x</span></html>`, 0)
	a.NoError(err)
	a.Contains(out, `<span title="a > b">`)
}

func TestFormatHTMLDocStringPreservesPre(t *testing.T) {
	a := require.New(t)
	src := "<html>\n<pre>\n  keep   this\n    exactly\n</pre>\n</html>"
	out, err := formatHTMLDocString(src, 0)
	a.NoError(err)
	// The whitespace-sensitive <pre> content must be preserved verbatim.
	a.Contains(out, "<pre>\n  keep   this\n    exactly\n</pre>")
}

func TestFormatHTMLDocStringMalformedFailsWithClearError(t *testing.T) {
	a := require.New(t)

	// An unclosed block element reports which tag(s) were left open.
	_, err := formatHTMLDocString("<html><p>unclosed paragraph</html>", 0)
	a.Error(err)
	a.Contains(err.Error(), "mismatched closing tag </html>")

	// A never-closed document lists the unclosed tags.
	_, err = formatHTMLDocString("<html><p>never closed", 0)
	a.Error(err)
	a.Contains(err.Error(), "unclosed tag(s)")

	// A mismatched inline end tag is reported clearly.
	_, err = formatHTMLDocString("<html><b>mismatched</i></html>", 0)
	a.Error(err)
	a.Contains(err.Error(), "mismatched closing tag </i>")

	// A stray closing tag with no matching opening tag is reported clearly.
	_, err = formatHTMLDocString("<div></div></div>", 0)
	a.Error(err)
	a.Contains(err.Error(), "unexpected closing tag </div>")

	// An unterminated preserved element (<pre>) is reported clearly.
	_, err = formatHTMLDocString("<html><pre>never closed", 0)
	a.Error(err)
	a.Contains(err.Error(), "unterminated <pre> element")
}

func TestFormatHTMLDocStringInlineElementsStayOnOneLine(t *testing.T) {
	a := require.New(t)

	// Inline elements and their surrounding text should be condensed onto one
	// line inside their block parent, rather than exploded one tag per line.
	src := "<html>\n<p>\nHello &amp; welcome to\n<b>Example</b>\n.\n</p>\n</html>"
	out, err := formatHTMLDocString(src, 0)
	a.NoError(err)
	a.Equal("<html>\n  <p>Hello &amp; welcome to <b>Example</b> .</p>\n</html>", out)

	// A link with an attribute and a trailing sentence collapses to a single line.
	src = "<html>\n<p>\nSee\n<a href=\"x\">Example</a>\nfor details.\n</p>\n</html>"
	out, err = formatHTMLDocString(src, 0)
	a.NoError(err)
	a.Equal("<html>\n  <p>See <a href=\"x\">Example</a> for details.</p>\n</html>", out)

	// A void inline element (<br/>) joins the surrounding text on the same line.
	src = "<html>\n<ul>\n<li>First:<br/>second.</li>\n</ul>\n</html>"
	out, err = formatHTMLDocString(src, 0)
	a.NoError(err)
	a.Equal("<html>\n  <ul>\n    <li>\n      First:<br/>second.\n    </li>\n  </ul>\n</html>", out)

	// Plain h4 headings compact like paragraphs.
	src = "<html>\n<h4>\nReference\n</h4>\n</html>"
	out, err = formatHTMLDocString(src, 0)
	a.NoError(err)
	a.Equal("<html>\n  <h4>Reference</h4>\n</html>", out)

	// Attributed h4 headings keep the block-style multiline layout.
	src = `<html><h4 class="ref">Reference</h4></html>`
	out, err = formatHTMLDocString(src, 0)
	a.NoError(err)
	a.Equal("<html>\n  <h4 class=\"ref\">\n    Reference\n  </h4>\n</html>", out)
}

func TestMaybeFormatHTMLString(t *testing.T) {
	a := require.New(t)

	// Non-HTML string: returned unchanged, no error.
	orig := `"just plain text"`
	out, ok, err := maybeFormatHTMLString(orig, 0)
	a.False(ok)
	a.NoError(err)
	a.Equal(orig, out)

	// HTML string with escaped quotes: formatted and re-escaped losslessly. The
	// inline <a> element stays on one line with its text.
	in := `"<html><a href=\"x\">y</a></html>"`
	out, ok, err = maybeFormatHTMLString(in, 0)
	a.True(ok)
	a.NoError(err)
	a.Equal("\"\n<html>\n  <a href=\\\"x\\\">y</a>\n</html>\"", out)

	// Malformed HTML string: returned unchanged, with a clear error describing
	// the problem so the caller can report it instead of silently emitting it.
	bad := `"<html><p>no close"`
	out, ok, err = maybeFormatHTMLString(bad, 0)
	a.False(ok)
	a.Error(err)
	a.Contains(err.Error(), "unclosed tag(s)")
	a.Equal(bad, out)
}
