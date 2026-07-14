// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package format

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// revisionListBlankLineRe matches the stray blank line(s) that are commonly left
// between the closing `</ul>` of a revisions list and the closing `</html>` of a
// Modelica annotation docstring. The trailing indentation before `</html>` is
// captured so it can be preserved. HTML tag names are matched case-insensitively.
var revisionListBlankLineRe = regexp.MustCompile(`(?i)(</ul>)[ \t]*\r?\n(?:[ \t]*\r?\n)+([ \t]*)(</html>)`)

// collapseRevisionListBlankLine removes the extra empty line(s) frequently left
// between the final `</ul>` and `</html>` of a revisions docstring, turning
// `</ul>\n\n</html>` into `</ul>\n</html>` (see issue #26). Content other than
// that blank-line run is left untouched (including the original tag casing and
// the indentation preceding `</html>`).
func collapseRevisionListBlankLine(s string) string {
	return revisionListBlankLineRe.ReplaceAllString(s, "${1}\n${2}${3}")
}

// voidElements are HTML elements that never have content or an end tag.
var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

// preserveElements are elements whose inner content is whitespace-sensitive and
// must be emitted verbatim (never reflowed or reindented).
var preserveElements = map[string]bool{
	"pre": true, "textarea": true, "script": true, "style": true,
}

// inlineElements are HTML phrasing/inline elements that are kept on the same line
// as the surrounding text rather than broken onto their own lines. This keeps
// short markup such as <b>bold</b> or <a href="...">link</a> condensed instead of
// exploding every tag onto a separate line.
var inlineElements = map[string]bool{
	"a": true, "abbr": true, "b": true, "bdi": true, "bdo": true, "br": true,
	"cite": true, "code": true, "data": true, "dfn": true, "em": true, "i": true,
	"img": true, "kbd": true, "label": true, "mark": true, "q": true, "rp": true,
	"rt": true, "ruby": true, "s": true, "samp": true, "small": true, "span": true,
	"strong": true, "sub": true, "sup": true, "time": true, "u": true, "var": true,
	"wbr": true, "tt": true, "big": true, "strike": true, "font": true,
	"acronym": true,
}

// compactBlockElements are block-level elements whose inline-only content should
// stay on the same line as the opening and closing tags.
var compactBlockElements = map[string]bool{
	"p": true,
}

// compactBlockElementsWithoutAttributes are compacted only when the opening tag
// has no attributes, e.g., <h4>Reference</h4> but not <h4 class="...">.
var compactBlockElementsWithoutAttributes = map[string]bool{
	"h4": true,
}

type htmlStackEntry struct {
	name          string
	inline        bool
	lineIndex     int
	hasAttributes bool
}

// isHTMLContent returns true if the (already unescaped) string content should be
// treated as HTML. Per Modelica convention this is when the first non-whitespace
// content is an opening <html> tag.
func isHTMLContent(decoded string) bool {
	trimmed := strings.TrimLeft(decoded, " \t\r\n\f\v")
	return len(trimmed) >= 5 && strings.EqualFold(trimmed[:5], "<html")
}

// unescapeModelicaString converts the raw content of a Modelica string literal
// (without its surrounding quotes) into the actual character data it represents,
// decoding the standard Modelica/C escape sequences.
func unescapeModelicaString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '"':
				b.WriteByte('"')
			case '\\':
				b.WriteByte('\\')
			case '\'':
				b.WriteByte('\'')
			case '?':
				b.WriteByte('?')
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			case 'a':
				b.WriteByte('\a')
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case 'v':
				b.WriteByte('\v')
			default:
				// unknown escape: keep both characters verbatim
				b.WriteByte('\\')
				b.WriteByte(s[i+1])
			}
			i++
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

// escapeModelicaString re-encodes character data into the content of a Modelica
// string literal (without surrounding quotes). Only backslash and double-quote
// require escaping; newlines and other whitespace are emitted literally, which
// matches how HTML docstrings are conventionally written.
func escapeModelicaString(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// htmlTagName extracts the lowercased tag name from a raw tag token such as
// "<div class=...>", "</div>" or "<br/>".
func htmlTagName(raw string) string {
	s := strings.TrimPrefix(raw, "<")
	s = strings.TrimPrefix(s, "/")
	end := len(s)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\f', '\v', '/', '>':
			end = i
			i = len(s)
		}
	}
	return strings.ToLower(s[:end])
}

// isSelfClosingTag returns true for a start tag written with a trailing slash
// (e.g., "<br/>").
func isSelfClosingTag(raw string) bool {
	t := strings.TrimRight(raw, " \t\r\n\f\v")
	return strings.HasSuffix(t, "/>")
}

// htmlStartTagHasAttributes reports whether a raw start tag has content after
// the tag name other than optional whitespace, slash, and closing angle bracket.
func htmlStartTagHasAttributes(raw string) bool {
	s := strings.TrimSpace(raw)
	s = strings.TrimPrefix(s, "<")
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\f', '\v', '/', '>':
			rest := strings.TrimSpace(s[i:])
			rest = strings.TrimSuffix(rest, ">")
			rest = strings.TrimSpace(rest)
			rest = strings.TrimSuffix(rest, "/")
			rest = strings.TrimSpace(rest)
			return rest != ""
		}
	}
	return false
}

// collapseInlineWhitespace collapses runs of ASCII whitespace to a single space
// while preserving a single leading/trailing space when the original text had one,
// so inline text joins cleanly with adjacent inline elements. Non-ASCII bytes
// (including UTF-8 NBSP and entity text like "&nbsp;") are preserved exactly so no
// content is lost. Boundary spaces are trimmed later when the inline run is flushed.
func collapseInlineWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			inSpace = true
		default:
			if inSpace {
				b.WriteByte(' ')
				inSpace = false
			}
			b.WriteByte(s[i])
		}
	}
	if inSpace {
		b.WriteByte(' ')
	}
	return b.String()
}

// formatHTMLDocString pretty-prints HTML content. Block-level elements, comments
// and whitespace-sensitive elements are placed on their own lines indented to their
// nesting depth (starting at baseIndent). Inline/phrasing elements (see
// inlineElements) and the text around them are kept together on a single line, so
// short markup such as <b>bold</b> or <a href="...">link</a> stays condensed.
// Tag/attribute/entity/text bytes are preserved exactly; only whitespace changes.
// Whitespace-sensitive elements (pre/textarea/script/style) are emitted verbatim.
//
// It returns a non-nil error (and the caller should leave the original untouched)
// when the HTML is malformed or unbalanced, so a bad docstring is never corrupted.
// The error explains what was wrong (e.g., an unexpected closing tag or an unclosed
// element) so the failure can be reported clearly.
func formatHTMLDocString(decoded string, baseIndent int) (string, error) {
	z := html.NewTokenizer(strings.NewReader(decoded))
	var lines []string
	var stack []htmlStackEntry
	depth := baseIndent

	var inline strings.Builder
	indentAt := func(d int) string { return strings.Repeat(spaceIndent, d) }

	// flushInline emits any accumulated inline run as its own line (trimmed of the
	// boundary spaces preserved during collapsing) and resets the buffer.
	flushInline := func() {
		text := strings.Trim(inline.String(), " \t\n\r\f\v")
		inline.Reset()
		if text != "" {
			lines = append(lines, indentAt(depth)+text)
		}
	}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				break
			}
			return "", fmt.Errorf("HTML tokenizer error: %v", z.Err())
		}
		raw := string(z.Raw())

		switch tt {
		case html.TextToken:
			inline.WriteString(collapseInlineWhitespace(raw))
		case html.SelfClosingTagToken:
			if inlineElements[htmlTagName(raw)] {
				inline.WriteString(raw)
			} else {
				flushInline()
				lines = append(lines, indentAt(depth)+raw)
			}
		case html.StartTagToken:
			name := htmlTagName(raw)
			if preserveElements[name] {
				flushInline()
				block, err := capturePreserved(z, name, raw)
				if err != nil {
					return "", err
				}
				lines = append(lines, indentAt(depth)+block)
				break
			}
			if inlineElements[name] {
				inline.WriteString(raw)
				if !voidElements[name] && !isSelfClosingTag(raw) {
					stack = append(stack, htmlStackEntry{name: name, inline: true})
				}
				break
			}
			// block element
			flushInline()
			lines = append(lines, indentAt(depth)+raw)
			if voidElements[name] || isSelfClosingTag(raw) {
				break
			}
			stack = append(stack, htmlStackEntry{
				name:          name,
				lineIndex:     len(lines) - 1,
				hasAttributes: htmlStartTagHasAttributes(raw),
			})
			depth++
		case html.EndTagToken:
			name := htmlTagName(raw)
			if voidElements[name] {
				// stray end tag for a void element: emit but don't dedent
				if inlineElements[name] {
					inline.WriteString(raw)
				} else {
					flushInline()
					lines = append(lines, indentAt(depth)+raw)
				}
				break
			}
			if len(stack) == 0 {
				return "", fmt.Errorf("unexpected closing tag </%s> (no matching opening tag)", name)
			}
			entry := stack[len(stack)-1]
			if entry.name != name {
				return "", fmt.Errorf("mismatched closing tag </%s> (expected </%s>)", name, entry.name)
			}
			stack = stack[:len(stack)-1]
			if entry.inline {
				inline.WriteString(raw)
				break
			}
			text := strings.Trim(inline.String(), " \t\n\r\f\v")
			canCompact := compactBlockElements[name] ||
				(compactBlockElementsWithoutAttributes[name] && !entry.hasAttributes)
			if canCompact && text != "" && entry.lineIndex == len(lines)-1 {
				inline.Reset()
				depth--
				lines[entry.lineIndex] += text + raw
				break
			}
			// block element
			flushInline()
			depth--
			lines = append(lines, indentAt(depth)+raw)
		case html.CommentToken, html.DoctypeToken:
			flushInline()
			lines = append(lines, indentAt(depth)+strings.TrimSpace(raw))
		}
	}

	flushInline()
	if len(stack) != 0 {
		names := make([]string, 0, len(stack))
		for _, entry := range stack {
			names = append(names, entry.name)
		}
		return "", fmt.Errorf("unclosed tag(s): <%s>", strings.Join(names, ">, <"))
	}
	return strings.Join(lines, "\n"), nil
}

// capturePreserved consumes tokens verbatim from the tokenizer until the end tag
// matching name is reached, returning the exact original bytes of the whole element
// (start tag + inner content + end tag). This keeps whitespace-sensitive content
// such as <pre> intact. It returns an error if the element is never closed.
func capturePreserved(z *html.Tokenizer, name, startRaw string) (string, error) {
	var b strings.Builder
	b.WriteString(startRaw)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return "", fmt.Errorf("unterminated <%s> element", name)
		}
		raw := string(z.Raw())
		b.WriteString(raw)
		if tt == html.EndTagToken && htmlTagName(raw) == name {
			return b.String(), nil
		}
	}
}

// maybeFormatHTMLString inspects a raw Modelica STRING token (including its
// surrounding quotes). Its return values are:
//   - (formatted, true, nil)  when the content is HTML and was reformatted;
//   - (tokenText, false, nil) when the content is not an HTML docstring (left as-is);
//   - (tokenText, false, err) when the content looks like HTML (starts with <html>)
//     but is malformed/unbalanced. The original token is returned unchanged and the
//     error describes the problem so the caller can report it clearly instead of
//     silently emitting unformatted HTML.
func maybeFormatHTMLString(tokenText string, baseIndent int) (string, bool, error) {
	if len(tokenText) < 2 || tokenText[0] != '"' || tokenText[len(tokenText)-1] != '"' {
		return tokenText, false, nil
	}
	inner := tokenText[1 : len(tokenText)-1]
	decoded := unescapeModelicaString(inner)
	if !isHTMLContent(decoded) {
		return tokenText, false, nil
	}
	formatted, err := formatHTMLDocString(decoded, baseIndent)
	if err != nil {
		return tokenText, false, err
	}
	return "\"\n" + escapeModelicaString(formatted) + "\"", true, nil
}
