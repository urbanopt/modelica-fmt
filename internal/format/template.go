// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package format

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Template dialects supported for .mot (templated Modelica) files.
const (
	// DialectJinja is the Jinja template dialect used by geojson-modelica-translator.
	DialectJinja = "jinja"
)

// Jinja templating constructs are not valid Modelica, so a .mot file cannot be
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
	jinjaControlRegex    = regexp.MustCompile(`{%.*?%}`)
	jinjaExpressionRegex = regexp.MustCompile(`{{.*?}}`)
	jinjaRawBlockRegex   = regexp.MustCompile(`(?s){% raw %}.*?{% endraw %}`)
	// The reverse regexes match three-or-more digits so that files with 1000+
	// substitutions (whose ids widen past the %03d minimum) are still restored
	// correctly, rather than leaving a stray trailing digit.
	commentedSubRegex = regexp.MustCompile(`/\*(JINJA_SUB_\d{3,})\*/`)
	normalSubRegex    = regexp.MustCompile(`JINJA_SUB_\d{3,}`)
)

// isTemplateFile reports whether the given path is a templated Modelica file.
func isTemplateFile(name string) bool {
	return strings.HasSuffix(name, ".mot")
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
// placeholders so the surrounding text remains valid Modelica.
func subControl(text string, sub *subMap) string {
	return jinjaControlRegex.ReplaceAllStringFunc(text, func(match string) string {
		return "/*" + sub.addSub(match) + "*/"
	})
}

// subExpression replaces Jinja expressions ({{ ... }}) with a bare placeholder
// identifier.
func subExpression(text string, sub *subMap) string {
	return jinjaExpressionRegex.ReplaceAllStringFunc(text, func(match string) string {
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
		seg = subControl(seg, sub)
		seg = subExpression(seg, sub)
		b.WriteString(seg)

		// The raw block itself: substitute only the control tags (including the
		// {% raw %}/{% endraw %} tags), leaving any {{ ... }} as literal text.
		raw := text[start:end]
		raw = subControl(raw, sub)
		b.WriteString(raw)

		prevEnd = end
	}

	// Remaining text after the last raw block.
	seg := text[prevEnd:]
	seg = subControl(seg, sub)
	seg = subExpression(seg, sub)
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
	text = commentedSubRegex.ReplaceAllString(text, "${1}")

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

// processTemplate formats a templated Modelica (.mot) file by running the
// substitute -> format -> reverse round trip described above. If the substituted
// text still cannot be parsed cleanly, the underlying formatter returns an error
// and the original file is left unchanged by the caller.
func processTemplate(text string, out io.Writer, config Config, dialect, filename string) error {
	sub := newSubMap()
	substituted, err := substituteTemplate(dialect, text, sub)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := formatModelica(substituted, &buf, config, filename); err != nil {
		return err
	}

	restored, err := reverseSub(buf.String(), sub)
	if err != nil {
		return err
	}

	_, err = io.WriteString(out, restored)
	return err
}
