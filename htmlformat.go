// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

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
// (e.g. "<br/>").
func isSelfClosingTag(raw string) bool {
	t := strings.TrimRight(raw, " \t\r\n\f\v")
	return strings.HasSuffix(t, "/>")
}

// collapseHTMLWhitespace trims and collapses runs of ASCII whitespace to a single
// space. Non-ASCII-whitespace bytes (including UTF-8 NBSP and entity text like
// "&nbsp;") are preserved exactly so no content is lost.
func collapseHTMLWhitespace(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	inSpace := false
	started := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			inSpace = true
		default:
			if inSpace && started {
				b.WriteByte(' ')
			}
			b.WriteByte(s[i])
			started = true
			inSpace = false
		}
	}
	return b.String()
}

// formatHTMLDocString pretty-prints HTML content, placing each tag, text node and
// comment on its own line indented to its nesting depth (starting at baseIndent).
// Tag/attribute/entity/text bytes are preserved exactly; only whitespace changes.
// Whitespace-sensitive elements (pre/textarea/script/style) are emitted verbatim.
//
// It returns ok=false (and the caller should leave the original untouched) when the
// HTML is malformed or unbalanced, so a bad docstring is never corrupted.
func formatHTMLDocString(decoded string, baseIndent int) (string, bool) {
	z := html.NewTokenizer(strings.NewReader(decoded))
	var lines []string
	var stack []string
	depth := baseIndent

	indentAt := func(d int) string { return strings.Repeat(spaceIndent, d) }

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				break
			}
			return "", false
		}
		raw := string(z.Raw())

		switch tt {
		case html.TextToken:
			text := collapseHTMLWhitespace(raw)
			if text == "" {
				continue
			}
			lines = append(lines, indentAt(depth)+text)
		case html.SelfClosingTagToken:
			lines = append(lines, indentAt(depth)+raw)
		case html.StartTagToken:
			name := htmlTagName(raw)
			if preserveElements[name] {
				block, ok := capturePreserved(z, name, raw)
				if !ok {
					return "", false
				}
				lines = append(lines, indentAt(depth)+block)
				break
			}
			lines = append(lines, indentAt(depth)+raw)
			if voidElements[name] || isSelfClosingTag(raw) {
				break
			}
			stack = append(stack, name)
			depth++
		case html.EndTagToken:
			name := htmlTagName(raw)
			if voidElements[name] {
				// stray end tag for a void element: emit but don't dedent
				lines = append(lines, indentAt(depth)+raw)
				break
			}
			if len(stack) == 0 || stack[len(stack)-1] != name {
				return "", false
			}
			stack = stack[:len(stack)-1]
			depth--
			lines = append(lines, indentAt(depth)+raw)
		case html.CommentToken, html.DoctypeToken:
			lines = append(lines, indentAt(depth)+strings.TrimSpace(raw))
		}
	}

	if len(stack) != 0 {
		return "", false
	}
	return strings.Join(lines, "\n"), true
}

// capturePreserved consumes tokens verbatim from the tokenizer until the end tag
// matching name is reached, returning the exact original bytes of the whole element
// (start tag + inner content + end tag). This keeps whitespace-sensitive content
// such as <pre> intact.
func capturePreserved(z *html.Tokenizer, name, startRaw string) (string, bool) {
	var b strings.Builder
	b.WriteString(startRaw)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			return "", false
		}
		raw := string(z.Raw())
		b.WriteString(raw)
		if tt == html.EndTagToken && htmlTagName(raw) == name {
			return b.String(), true
		}
	}
}

// maybeFormatHTMLString inspects a raw Modelica STRING token (including its
// surrounding quotes). If its content is HTML, it returns a reformatted string
// literal and ok=true; otherwise it returns the original token unchanged.
func maybeFormatHTMLString(tokenText string, baseIndent int) (string, bool) {
	if len(tokenText) < 2 || tokenText[0] != '"' || tokenText[len(tokenText)-1] != '"' {
		return tokenText, false
	}
	inner := tokenText[1 : len(tokenText)-1]
	decoded := unescapeModelicaString(inner)
	if !isHTMLContent(decoded) {
		return tokenText, false
	}
	formatted, ok := formatHTMLDocString(decoded, baseIndent)
	if !ok {
		return tokenText, false
	}
	return "\"\n" + escapeModelicaString(formatted) + "\"", true
}
