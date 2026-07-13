// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package format

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

// Template dialects supported for .mot/.mopt (templated Modelica) files.
const (
	// DialectJinja is the Jinja template dialect used by geojson-modelica-translator.
	DialectJinja = "jinja"
)

// Jinja templating constructs are not valid Modelica, so a template file cannot be
// parsed/formatted directly. The strategy (ported from geojson-modelica-translator's
// management/format_modelica_files.py) is a substitute -> format -> reverse round trip:
//
//  1. Replace every Jinja construct with a unique placeholder identifier that the
//     Modelica lexer accepts. Control statements ({% ... %}) become commented-out
//     placeholders (/*JINJA_SUB_NNN*/) because control flow isn't valid Modelica;
//     expressions ({{ ... }}) become a bare identifier (JINJA_SUB_NNN). Inside a
//     {% raw %} ... {% endraw %} block only the control tags are substituted, since
//     {{ ... }} there is literal text.
//  2. Run the normal Modelica formatter on the now-parseable text.
//  3. Reverse the substitutions, stripping the /* */ wrappers around control
//     placeholders and restoring each placeholder's original text.
var (
	jinjaControlRegex         = regexp.MustCompile(`{%.*?%}`)
	jinjaExpressionRegex      = regexp.MustCompile(`{{.*?}}`)
	dollarExpressionRegex     = regexp.MustCompile(`\$\{.*?}`)
	jinjaExpressionLineRegex  = regexp.MustCompile(`^{{.*?}}$`)
	dollarExpressionLineRegex = regexp.MustCompile(`^\$\{.*?}$`)
	jinjaRawBlockRegex        = regexp.MustCompile(`(?s){% raw %}.*?{% endraw %}`)
	// Jinja loops often emit comma-separated Modelica arrays with
	// {% if not loop.last %},{% endif %}. If only the control tags are
	// substituted, the comma remains and produces an invalid trailing-comma
	// array in the temporary Modelica text.
	jinjaLoopCommaRegex = regexp.MustCompile(`{%\s*if\s+not\s+loop\.last\s*%}\s*,\s*{%\s*endif\s*%}`)
	// The reverse regexes match three-or-more digits so that files with 1000+
	// substitutions (whose ids widen past the %03d minimum) are still restored
	// correctly, rather than leaving a stray trailing digit.
	standaloneLineCommentedSubRegex = regexp.MustCompile(`(?m)(^[ \t]*)//\s*(JINJA_SUB_\d{3,})`)
	inlineLineCommentedSubRegex     = regexp.MustCompile(`(^|[^:])//\s*(JINJA_SUB_\d{3,})`)
	blockCommentedSubRegex          = regexp.MustCompile(`/\*(JINJA_SUB_\d{3,})\*/`)
	normalSubRegex                  = regexp.MustCompile(`JINJA_SUB_\d{3,}`)
)

// isTemplateFile reports whether the given path is a templated Modelica file.
func isTemplateFile(name string) bool {
	return strings.HasSuffix(name, ".mot") || strings.HasSuffix(name, ".mopt")
}

// subMap manages the mapping between placeholder identifiers and the original
// template text they replaced.
type subMap struct {
	curID int
	m     map[string]string
}

func newSubMap() *subMap {
	return &subMap{curID: 1, m: make(map[string]string)}
}

// addSub registers a substitution and returns the placeholder identifier.
func (s *subMap) addSub(text string) string {
	id := fmt.Sprintf("JINJA_SUB_%03d", s.curID)
	s.m[id] = text
	s.curID++
	return id
}

// getText returns the original text for a placeholder identifier.
func (s *subMap) getText(id string) (string, error) {
	text, ok := s.m[id]
	if !ok {
		return "", fmt.Errorf(
			"placeholder %q was not found in the substitution map "+
				"(possibly a false-positive placeholder match in the source)", id)
	}
	return text, nil
}

// subControl replaces Jinja control statements ({% ... %}) with commented-out
// placeholders so the surrounding text remains valid Modelica. Standalone
// controls use line comments so the formatter keeps them separated from
// adjacent Modelica comments and declarations; inline controls use block
// comments so they do not swallow the rest of the line.
func subControl(text string, sub *subMap) string {
	return subControlSegment(text, text, 0, sub)
}

func subControlSegment(fullText, text string, offset int, sub *subMap) string {
	var b strings.Builder
	lineOffset := offset
	for len(text) > 0 {
		line := text
		rest := ""
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			line = text[:i+1]
			rest = text[i+1:]
		}
		b.WriteString(subControlInLine(fullText, line, lineOffset, sub))
		lineOffset += len(line)
		text = rest
	}
	return b.String()
}

type controlSpan struct {
	start       int
	end         int
	lineComment bool
}

func subControlInLine(fullText, line string, lineOffset int, sub *subMap) string {
	var spans []controlSpan
	for _, span := range jinjaLoopCommaRegex.FindAllStringIndex(line, -1) {
		spans = append(spans, controlSpan{start: span[0], end: span[1]})
	}

	for _, span := range jinjaControlRegex.FindAllStringIndex(line, -1) {
		insideLoopComma := false
		for _, loopSpan := range spans {
			if span[0] >= loopSpan.start && span[1] <= loopSpan.end {
				insideLoopComma = true
				break
			}
		}
		if insideLoopComma {
			continue
		}
		absStart := lineOffset + span[0]
		absEnd := lineOffset + span[1]
		spans = append(spans, controlSpan{
			start:       span[0],
			end:         span[1],
			lineComment: isStandaloneConstructAt(fullText, absStart, absEnd),
		})
	}

	if len(spans) == 0 {
		return line
	}

	sort.Slice(spans, func(i, j int) bool {
		return spans[i].start < spans[j].start
	})

	var b strings.Builder
	prevEnd := 0
	for _, span := range spans {
		b.WriteString(line[prevEnd:span.start])
		b.WriteString(placeholderComment(line[span.start:span.end], sub, span.lineComment))
		prevEnd = span.end
	}
	b.WriteString(line[prevEnd:])
	return b.String()
}

func placeholderComment(text string, sub *subMap, lineComment bool) string {
	if lineComment {
		return "// " + sub.addSub(text)
	}
	return "/*" + sub.addSub(text) + "*/"
}

func isStandaloneConstructAt(text string, start, end int) bool {
	lineStart := strings.LastIndexByte(text[:start], '\n') + 1
	lineEnd := len(text)
	if i := strings.IndexByte(text[end:], '\n'); i >= 0 {
		lineEnd = end + i
	}
	before := strings.TrimSpace(text[lineStart:start])
	after := strings.TrimSpace(text[end:lineEnd])
	if after != "" {
		return false
	}
	return before == "" || strings.HasSuffix(before, "{")
}

func subRawBlockControls(fullText, rawBlock string, blockStart int, sub *subMap) string {
	var b strings.Builder
	lineOffset := blockStart
	for len(rawBlock) > 0 {
		line := rawBlock
		rest := ""
		if i := strings.IndexByte(rawBlock, '\n'); i >= 0 {
			line = rawBlock[:i+1]
			rest = rawBlock[i+1:]
		}
		b.WriteString(subControlInLine(fullText, line, lineOffset, sub))
		lineOffset += len(line)
		rawBlock = rest
	}
	return b.String()
}

// subExpressions replaces template expressions with placeholders. GMT snippet
// fields that occupy a whole line are commented out because they expand to
// generated declarations or equations. Everywhere else expressions become bare
// identifiers so they can stand in for names, scalar values, and string
// contents.
func subExpressions(text string, sub *subMap) string {
	var b strings.Builder
	for len(text) > 0 {
		line := text
		rest := ""
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			line = text[:i+1]
			rest = text[i+1:]
		}
		b.WriteString(subExpressionsInLine(line, sub))
		text = rest
	}
	return b.String()
}

func subExpressionsInLine(line string, sub *subMap) string {
	newline := ""
	body := line
	if strings.HasSuffix(body, "\n") {
		newline = "\n"
		body = strings.TrimSuffix(body, "\n")
	}
	if strings.HasSuffix(body, "\r") {
		newline = "\r" + newline
		body = strings.TrimSuffix(body, "\r")
	}

	prefixLen := len(body) - len(strings.TrimLeft(body, " \t"))
	suffixLen := len(body) - len(strings.TrimRight(body, " \t"))
	if prefixLen+suffixLen > len(body) {
		return line
	}
	prefix := body[:prefixLen]
	suffix := body[len(body)-suffixLen:]
	trimmed := body[prefixLen : len(body)-suffixLen]

	if isGeneratedSnippetExpression(trimmed) {
		return prefix + placeholderComment(trimmed, sub, true) + suffix + newline
	}

	return subExpression(body, sub) + newline
}

func isGeneratedSnippetExpression(text string) bool {
	if !jinjaExpressionLineRegex.MatchString(text) && !dollarExpressionLineRegex.MatchString(text) {
		return false
	}
	return strings.Contains(text, ".instance") ||
		strings.Contains(text, ".component_definitions") ||
		strings.Contains(text, ".connect_statements")
}

// subExpression replaces inline template expressions with a bare placeholder
// identifier.
func subExpression(text string, sub *subMap) string {
	text = jinjaExpressionRegex.ReplaceAllStringFunc(text, func(match string) string {
		return sub.addSub(match)
	})
	return dollarExpressionRegex.ReplaceAllStringFunc(text, func(match string) string {
		return sub.addSub(match)
	})
}

// substituteJinja replaces all Jinja constructs in text with placeholders,
// treating {% raw %} ... {% endraw %} blocks specially (only their control tags
// are substituted; {{ ... }} inside a raw block is literal and left untouched).
func substituteJinja(text string, sub *subMap) string {
	var b strings.Builder
	prevEnd := 0
	for _, span := range jinjaRawBlockRegex.FindAllStringIndex(text, -1) {
		start, end := span[0], span[1]

		// Text before the raw block: substitute both control and expressions.
		seg := text[prevEnd:start]
		seg = subControlSegment(text, seg, prevEnd, sub)
		seg = subExpressions(seg, sub)
		b.WriteString(seg)

		// The raw block itself: substitute only the control tags (including the
		// {% raw %}/{% endraw %} tags), leaving any {{ ... }} as literal text.
		raw := text[start:end]
		raw = subRawBlockControls(text, raw, start, sub)
		b.WriteString(raw)

		prevEnd = end
	}

	// Remaining text after the last raw block.
	seg := text[prevEnd:]
	seg = subControl(seg, sub)
	seg = subExpressions(seg, sub)
	b.WriteString(seg)

	return b.String()
}

// substituteTemplate applies the substitution strategy for the requested dialect.
func substituteTemplate(dialect, text string, sub *subMap) (string, error) {
	switch dialect {
	case DialectJinja, "":
		return substituteJinja(text, sub), nil
	default:
		return "", fmt.Errorf("unsupported template dialect %q (supported: jinja)", dialect)
	}
}

// reverseSub reverses the substitutions: it strips the /* */ wrappers around
// control-statement placeholders and then replaces every placeholder with its
// original template text.
func reverseSub(text string, sub *subMap) (string, error) {
	// Remove the comment wrappers around control-statement placeholders so that
	// both commented and bare placeholders are restored uniformly below.
	text = standaloneLineCommentedSubRegex.ReplaceAllString(text, "${1}${2}")
	text = inlineLineCommentedSubRegex.ReplaceAllString(text, "${1} ${2}")
	text = blockCommentedSubRegex.ReplaceAllString(text, "${1}")
	text = trimWhitespaceBeforeRawPunctuation(text, sub)

	var firstErr error
	restored := normalSubRegex.ReplaceAllStringFunc(text, func(match string) string {
		orig, err := sub.getText(match)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			return match
		}
		return orig
	})
	if firstErr != nil {
		return "", firstErr
	}
	return restored, nil
}

func trimWhitespaceBeforeRawPunctuation(text string, sub *subMap) string {
	matches := normalSubRegex.FindAllStringIndex(text, -1)
	if len(matches) < 2 {
		return text
	}

	var b strings.Builder
	last := 0
	for i := 0; i < len(matches)-1; i++ {
		cur := matches[i]
		next := matches[i+1]
		between := text[cur[1]:next[0]]
		if strings.TrimSpace(between) != "" {
			continue
		}
		if !isRawControlPlaceholder(text[next[0]:next[1]], sub) {
			continue
		}
		if next[1] >= len(text) || !isNoSpaceBeforeByte(text[next[1]]) {
			continue
		}

		b.WriteString(text[last:cur[1]])
		last = next[0]
	}
	b.WriteString(text[last:])
	return b.String()
}

func isRawControlPlaceholder(id string, sub *subMap) bool {
	orig, ok := sub.m[id]
	return ok && strings.TrimSpace(orig) == "{% raw %}"
}

func isNoSpaceBeforeByte(b byte) bool {
	return b == ')'
}

func sameNonWhitespaceContent(a, b string) bool {
	ia, ib := 0, 0
	for {
		for ia < len(a) && isWhitespace(a[ia]) {
			ia++
		}
		for ib < len(b) && isWhitespace(b[ib]) {
			ib++
		}
		if ia == len(a) || ib == len(b) {
			return ia == len(a) && ib == len(b)
		}
		if a[ia] != b[ib] {
			return false
		}
		ia++
		ib++
	}
}

func isWhitespace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\f', '\v':
		return true
	default:
		return false
	}
}

// processTemplate formats a templated Modelica file by running the
// substitute -> format -> reverse round trip described above. If the substituted
// text still cannot be parsed cleanly, or if the round trip would change
// non-whitespace content, the original text is emitted unchanged so
// directory-wide formatting can continue across non-Modelica or unsupported
// templates without corrupting them.
func processTemplate(text string, out io.Writer, config Config, dialect, filename string) error {
	sub := newSubMap()
	substituted, err := substituteTemplate(dialect, text, sub)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := formatModelica(substituted, &buf, config, filename); err != nil {
		_, writeErr := io.WriteString(out, text)
		return writeErr
	}

	restored, err := reverseSub(buf.String(), sub)
	if err != nil {
		return err
	}

	if !sameNonWhitespaceContent(text, restored) {
		_, writeErr := io.WriteString(out, text)
		return writeErr
	}

	_, err = io.WriteString(out, restored)
	return err
}
