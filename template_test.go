package main

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// templateFileTests exercise the .mot (Jinja-templated Modelica) pipeline against
// golden outputs. These fixtures and their expected results mirror the behavior of
// geojson-modelica-translator's management/format_modelica_files.py so that GMT can
// eventually drop its wrapper and call modelicafmt directly.
var templateFileTests = []struct {
	sourceFile      string
	outFile         string
	formatterConfig Config
}{
	// simple {{ ... }} expressions
	{"gmt-boiler-polynomial.mot", "gmt-boiler-polynomial-out.mot", Config{-1, false}},
	// {% if %}/{% else %}/{% endif %} control statements plus expressions
	{"gmt-design-data-series.mot", "gmt-design-data-series-out.mot", Config{-1, false}},
	// {% raw %} ... {% endraw %} blocks interleaved with expressions
	{"gmt-cooling-indirect.mot", "gmt-cooling-indirect-out.mot", Config{-1, false}},
}

func TestFormattingTemplateExamples(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			testSourceFile := path.Join("examples", testCase.sourceFile)
			expectedOutFile := path.Join("examples", testCase.outFile)
			actualOutFile := path.Join(outputDir, testCase.outFile)
			file, err := os.Create(actualOutFile)
			a.NoError(err)
			defer file.Close()

			err = processFile(testSourceFile, file, testCase.formatterConfig)
			a.NoError(err)

			diff, err := diffFiles(expectedOutFile, actualOutFile)
			a.NoError(err)
			a.Len(diff, 0, "File diff should be empty")
		})
	}
}

// TestTemplateFormattingIsIdempotent ensures that formatting an already-formatted
// .mot file produces no further changes.
func TestTemplateFormattingIsIdempotent(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.outFile, func(t *testing.T) {
			formattedFile := path.Join("examples", testCase.outFile)
			var out bytes.Buffer
			err := processFile(formattedFile, &out, testCase.formatterConfig)
			a.NoError(err)

			expected, err := os.ReadFile(formattedFile)
			a.NoError(err)
			a.Equal(string(expected), out.String(), "re-formatting a formatted .mot should be a no-op")
		})
	}
}

// TestTemplatePreservesNonWhitespaceContent mirrors GMT's own invariant: after
// formatting a .mot file the only differences should be whitespace.
func TestTemplatePreservesNonWhitespaceContent(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			source, err := os.ReadFile(path.Join("examples", testCase.sourceFile))
			a.NoError(err)

			var out bytes.Buffer
			err = processFile(path.Join("examples", testCase.sourceFile), &out, testCase.formatterConfig)
			a.NoError(err)

			a.Equal(stripWhitespace(string(source)), stripWhitespace(out.String()),
				"original and formatted .mot should have identical non-whitespace content")
		})
	}
}

// TestTemplateUnsupportedDialectReturnsError ensures an unknown template dialect
// is rejected instead of silently corrupting the file.
func TestTemplateUnsupportedDialectReturnsError(t *testing.T) {
	a := require.New(t)
	sub := newSubMap()
	_, err := substituteTemplate("mako", "model M end M;", sub)
	a.Error(err)
	a.Contains(err.Error(), "mako")
}

// TestTemplateUnformattableReturnsError ensures a template that cannot be made
// parseable via preprocessing (matching GMT's SKIP_FILES) returns an error rather
// than emitting garbled output, so the caller leaves the file unchanged.
func TestTemplateUnformattableReturnsError(t *testing.T) {
	a := require.New(t)
	var out bytes.Buffer
	err := processFile(path.Join("examples", "gmt-district-energy-system.mot"), &out, Config{-1, false})
	a.Error(err, "known-unformattable .mot should return an error")
}

func TestSubstituteReverseRoundTrip(t *testing.T) {
	a := require.New(t)
	original := "within {{ project_name }}.Foo;\n{% if x %}Real a=1;{% else %}Real a=2;{% endif %}\n"

	sub := newSubMap()
	substituted, err := substituteTemplate(dialectJinja, original, sub)
	a.NoError(err)
	// Control statements are commented out; expressions become bare identifiers.
	a.NotContains(substituted, "{%")
	a.NotContains(substituted, "{{")
	a.Contains(substituted, "/*JINJA_SUB_")
	a.Contains(substituted, "JINJA_SUB_")

	restored, err := reverseSub(substituted, sub)
	a.NoError(err)
	a.Equal(original, restored, "substitution round trip should restore the original text")
}

// TestRawBlockLeavesExpressionsUntouched verifies that {{ ... }} inside a
// {% raw %} ... {% endraw %} block is treated as literal text (not substituted),
// while the surrounding control tags are.
func TestRawBlockLeavesExpressionsUntouched(t *testing.T) {
	a := require.New(t)
	original := "a={{ outside }}\n{% raw %}b={{ inside }}{% endraw %}\n"

	sub := newSubMap()
	substituted, err := substituteTemplate(dialectJinja, original, sub)
	a.NoError(err)

	// The literal expression inside the raw block must remain untouched.
	a.Contains(substituted, "b={{ inside }}")
	// The expression outside the raw block must be substituted.
	a.NotContains(substituted, "{{ outside }}")
	// The raw/endraw control tags must be substituted (commented out).
	a.NotContains(substituted, "{% raw %}")
	a.NotContains(substituted, "{% endraw %}")

	restored, err := reverseSub(substituted, sub)
	a.NoError(err)
	a.Equal(original, restored)
}

// TestReverseSubHandlesWidePlaceholders ensures placeholders wider than the
// %03d minimum (i.e. 1000+ substitutions) are restored correctly rather than
// leaving a stray trailing digit.
func TestReverseSubHandlesWidePlaceholders(t *testing.T) {
	a := require.New(t)
	sub := newSubMap()
	sub.m["JINJA_SUB_1000"] = "{{ wide_expr }}"
	sub.m["JINJA_SUB_1001"] = "{% if wide %}"

	restored, err := reverseSub("a=JINJA_SUB_1000;\n/*JINJA_SUB_1001*/\n", sub)
	a.NoError(err)
	a.Equal("a={{ wide_expr }};\n{% if wide %}\n", restored)
}

func stripWhitespace(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case ' ', '\t', '\n', '\r', '\f', '\v':
			continue
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
