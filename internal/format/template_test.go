package format

import (
	"bytes"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// templateFileTests exercise the Jinja-templated Modelica pipeline against golden
// outputs. These fixtures and their expected results mirror the behavior of
// geojson-modelica-translator's management/format_modelica_files.py so that GMT can
// eventually drop its wrapper and call modelicafmt directly.
var templateFileTests = []struct {
	sourceFile      string
	outFile         string
	formatterConfig Config
}{
	// simple {{ ... }} expressions
	{"gmt-boiler-polynomial.mot", "gmt-boiler-polynomial-out.mot", Config{-1, false, false, false}},
	// {% if %}/{% else %}/{% endif %} control statements plus expressions
	{"gmt-design-data-series.mot", "gmt-design-data-series-out.mot", Config{-1, false, false, false}},
	// {% raw %} ... {% endraw %} blocks interleaved with expressions
	{"gmt-cooling-indirect.mot", "gmt-cooling-indirect-out.mot", Config{-1, false, false, false}},
	// {% for %} loops + a filter expression + many {% raw %} blocks (SpawnBuilding)
	{"gmt-spawn-building.mot", "gmt-spawn-building-out.mot", Config{-1, false, false, false}},
	// {% if %} control plus a dozen {% raw %} blocks in a large file (TimeSeriesBuilding)
	{"gmt-time-series-building.mot", "gmt-time-series-building-out.mot", Config{-1, false, false, false}},
	// {% for %} loop over array elements with a conditional comma.
	{"gmt-dhc-5g-wh-ghx-hpdirectcooling-variable-dist.mot", "gmt-dhc-5g-wh-ghx-hpdirectcooling-variable-dist-out.mot", Config{-1, false, false, false}},
	{"gmt-hptrio-variable-dist.mot", "gmt-hptrio-variable-dist-out.mot", Config{-1, false, false, false}},
}

func TestFormattingTemplateExamples(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			testSourceFile := path.Join("testdata", testCase.sourceFile)
			expectedOutFile := path.Join("testdata", testCase.outFile)
			actualOutFile := path.Join(outputDir, testCase.outFile)
			file, err := os.Create(actualOutFile)
			a.NoError(err)
			defer file.Close()

			err = ProcessFile(testSourceFile, file, testCase.formatterConfig, DialectJinja)
			a.NoError(err)

			diff, err := diffFiles(expectedOutFile, actualOutFile)
			a.NoError(err)
			a.Len(diff, 0, "File diff should be empty")
		})
	}
}

// TestTemplateFormattingIsIdempotent ensures that formatting an already-formatted
// template file produces no further changes.
func TestTemplateFormattingIsIdempotent(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.outFile, func(t *testing.T) {
			formattedFile := path.Join("testdata", testCase.outFile)
			var out bytes.Buffer
			err := ProcessFile(formattedFile, &out, testCase.formatterConfig, DialectJinja)
			a.NoError(err)

			expected, err := os.ReadFile(formattedFile)
			a.NoError(err)
			a.Equal(string(expected), out.String(), "re-formatting a formatted .mot should be a no-op")
		})
	}
}

// TestTemplatePreservesNonWhitespaceContent mirrors GMT's own invariant: after
// formatting a template file the only differences should be whitespace.
func TestTemplatePreservesNonWhitespaceContent(t *testing.T) {
	a := require.New(t)
	for _, testCase := range templateFileTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			source, err := os.ReadFile(path.Join("testdata", testCase.sourceFile))
			a.NoError(err)

			var out bytes.Buffer
			err = ProcessFile(path.Join("testdata", testCase.sourceFile), &out, testCase.formatterConfig, DialectJinja)
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

// previouslySkippedTemplateTests are real GMT template files that used to fail
// directory-wide formatting. They should now process without returning an error:
// Modelica templates are formatted, and non-Modelica template scripts are passed
// through unchanged.
var previouslySkippedTemplateTests = []struct {
	sourceFile string
}{
	{"gmt-district-energy-system.mot"},
	{"gmt-run-spawn-building.mot"},
}

func TestPreviouslySkippedTemplatesDoNotError(t *testing.T) {
	a := require.New(t)
	for _, testCase := range previouslySkippedTemplateTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			source, err := os.ReadFile(path.Join("testdata", testCase.sourceFile))
			a.NoError(err)

			var out bytes.Buffer
			err = ProcessFile(path.Join("testdata", testCase.sourceFile), &out, Config{-1, false, false, false}, DialectJinja)
			a.NoError(err)
			a.Equal(stripWhitespace(string(source)), stripWhitespace(out.String()),
				"processing a .mot template must preserve non-whitespace content")
		})
	}
}

func TestStandaloneControlTagsKeepLineBoundaries(t *testing.T) {
	a := require.New(t)
	var out bytes.Buffer
	err := ProcessFile(
		path.Join("testdata", "gmt-district-energy-system.mot"),
		&out,
		Config{-1, false, false, false},
		DialectJinja)
	a.NoError(err)

	formatted := out.String()
	a.NotContains(formatted, "{% for model in models %}//")
	a.NotContains(formatted, "{{ model.instance }}//")
	a.NotContains(formatted, "{% endfor %}// Model dependencies")
	a.NotContains(formatted, "{% for coupling in couplings %}//")
	a.NotContains(formatted, "{{ coupling.component_definitions }}//")
	a.NotContains(formatted, "{% endfor %}equation")
	a.NotContains(formatted, "{{ coupling.connect_statements }} //")
	a.NotContains(formatted, "{% endfor %} annotation")

	a.Regexp(`(?s)\{% for model in models %\}\s*//`, formatted)
	a.Regexp(`(?s)\{\{ model\.instance \}\}\s*//`, formatted)
	a.Regexp(`(?s)\{% endfor %\}\s*// Model dependencies`, formatted)
	a.Regexp(`(?s)\{% endfor %\}\s*equation`, formatted)
}

func TestStandaloneRawTagsKeepLineBoundaries(t *testing.T) {
	a := require.New(t)
	original := `model RawBlock
{% raw %}
  RealInput u annotation (Placement(transformation(extent={{-240,-40},{-200,0}})));
{% endraw %}
end RawBlock;
`
	var out bytes.Buffer
	err := processTemplate(original, &out, Config{-1, false, false, false}, DialectJinja, "raw-block.mot")
	a.NoError(err)

	formatted := out.String()
	a.Contains(formatted, "{% raw %}\n")
	a.Contains(formatted, "\n{% endraw %}\n")
	a.NotContains(formatted, "{% raw %}RealInput")
	a.NotContains(formatted, "));{% endraw %}")
}

func TestMoptFilesUseTemplatePipeline(t *testing.T) {
	a := require.New(t)
	sourceFile := path.Join(outputDir, "template-extension.mopt")
	a.NoError(os.WriteFile(sourceFile, []byte("model {{ model_name }}\nend {{ model_name }};\n"), 0644))
	defer os.Remove(sourceFile)

	var out bytes.Buffer
	err := ProcessFile(sourceFile, &out, Config{-1, false, false, false}, DialectJinja)
	a.NoError(err)
	a.Equal("model {{ model_name }}\nend {{ model_name }};\n", out.String())
}

// TestNonModelicaTemplatePassesThrough specifically guards against silent data
// loss: a non-empty template file that is not actually Modelica must be emitted
// unchanged instead of being overwritten with an empty file.
func TestNonModelicaTemplatePassesThrough(t *testing.T) {
	a := require.New(t)
	source, err := os.ReadFile(path.Join("testdata", "gmt-run-spawn-building.mot"))
	a.NoError(err)

	var out bytes.Buffer
	err = ProcessFile(path.Join("testdata", "gmt-run-spawn-building.mot"), &out, Config{-1, false, false, false}, DialectJinja)
	a.NoError(err)
	a.Equal(string(source), out.String())
}

func TestTemplateWithMeaningfulDiffPassesThrough(t *testing.T) {
	a := require.New(t)
	source, err := os.ReadFile(path.Join("testdata", "gmt-district-energy-system.mot"))
	a.NoError(err)

	var out bytes.Buffer
	err = ProcessFile(path.Join("testdata", "gmt-district-energy-system.mot"), &out, Config{-1, false, false, false}, DialectJinja)
	a.NoError(err)
	a.Equal(string(source), out.String())
}

func TestSubstituteReverseRoundTrip(t *testing.T) {
	a := require.New(t)
	original := "within {{ project_name }}.Foo;\n{% if x %}Real a=1;{% else %}Real a=2;{% endif %}\n"

	sub := newSubMap()
	substituted, err := substituteTemplate(DialectJinja, original, sub)
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
	substituted, err := substituteTemplate(DialectJinja, original, sub)
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

func TestLoopConditionalCommaRoundTrip(t *testing.T) {
	a := require.New(t)
	original := `parameter String filNam[nBui]={
    {% for building in data["building_load_files"] %}
    "{{ building }}"{% if not loop.last %},{% endif %}
    {% endfor %}{% raw %}}
    "Library paths";{% endraw %}`

	sub := newSubMap()
	substituted, err := substituteTemplate(DialectJinja, original, sub)
	a.NoError(err)

	a.NotContains(substituted, ",", "conditional loop comma should not leak into temporary Modelica")
	a.Contains(substituted, `"JINJA_SUB_`)

	restored, err := reverseSub(substituted, sub)
	a.NoError(err)
	a.Equal(original, restored)
}

func TestLoopedArrayDoesNotCreateJinjaExpressionStart(t *testing.T) {
	a := require.New(t)
	original := `model LoopedArray
  parameter String filNam[nBui]={
    {% for building in data["building_load_files"] %}
    "{{ building }}"{% if not loop.last %},{% endif %}
    {% endfor %}{% raw %}}
    "Library paths";{% endraw %}
end LoopedArray;
`
	var out bytes.Buffer
	err := processTemplate(original, &out, Config{-1, false, false, false}, DialectJinja, "looped-array.mot")
	a.NoError(err)

	formatted := out.String()
	a.NotContains(formatted, "{{%", "a literal Modelica array opener must not merge with a Jinja control tag")
	a.Contains(formatted, "{% for building in data[\"building_load_files\"] %}")
}

// TestSingleElementArrayExpressionKeepsDisambiguatingSpace guards against the
// formatter collapsing `{ {{ expr }} }` into `{{{ expr }}}`. The Modelica
// array formatter normally strips the space inside a single-element array
// literal, but doing so here would make the restored Jinja expression's own
// `{{`/`}}` delimiters collide with the surrounding array braces, producing
// an ambiguous triple-brace sequence when the template is later rendered.
func TestSingleElementArrayExpressionKeepsDisambiguatingSpace(t *testing.T) {
	a := require.New(t)
	original := `model NominalArray
  parameter Real a[3] = fill({ {{ data["nominal_values"]["boiler_efficiency"] }} });
end NominalArray;
`
	var out bytes.Buffer
	err := processTemplate(original, &out, Config{-1, false, false, false}, DialectJinja, "nominal-array.mot")
	a.NoError(err)

	formatted := out.String()
	a.NotContains(formatted, `{{{`, "restored Jinja expression must not merge with a literal Modelica array brace")
	a.NotContains(formatted, `}}}`, "restored Jinja expression must not merge with a literal Modelica array brace")
	a.Contains(formatted, `{ {{ data["nominal_values"]["boiler_efficiency"] }} }`)
}

func TestInlineExpressionBeforeRawPunctuationDoesNotGainRenderedSpace(t *testing.T) {
	a := require.New(t)
	original := `model PumpTemplate
  Pump pumSto(
    dp_nominal={{ data["source_pump_dp_nominal"] }}{% raw %})
    "Bore field pump";
{% endraw %}end PumpTemplate;
`
	var out bytes.Buffer
	err := processTemplate(original, &out, Config{-1, false, false, false}, DialectJinja, "pump-template.mot")
	a.NoError(err)

	formatted := out.String()
	a.Contains(formatted, `{{ data["source_pump_dp_nominal"] }}{% raw %})`)
	a.NotContains(formatted, `{{ data["source_pump_dp_nominal"] }} {% raw %})`)
}

func TestGeneratedSnippetExpressionsAreCommented(t *testing.T) {
	a := require.New(t)
	original := "{{ model.instance }}\n{{ coupling.component_definitions }}\n{{ coupling.connect_statements }}\n{{ data['lCon'] }}\n"

	sub := newSubMap()
	substituted, err := substituteTemplate(DialectJinja, original, sub)
	a.NoError(err)

	a.Contains(substituted, "// JINJA_SUB_001")
	a.Contains(substituted, "// JINJA_SUB_002")
	a.Contains(substituted, "// JINJA_SUB_003")
	a.Contains(substituted, "\nJINJA_SUB_004\n")

	restored, err := reverseSub(substituted, sub)
	a.NoError(err)
	a.Equal(original, restored)
}

func TestDollarExpressionRoundTrip(t *testing.T) {
	a := require.New(t)
	original := "within ${project_name};\nmodel ${model_name}\n  Real x=${value};\nend ${model_name};\n"

	sub := newSubMap()
	substituted, err := substituteTemplate(DialectJinja, original, sub)
	a.NoError(err)
	a.NotContains(substituted, "${")
	a.Contains(substituted, "within JINJA_SUB_001;")

	restored, err := reverseSub(substituted, sub)
	a.NoError(err)
	a.Equal(original, restored)
}

// TestReverseSubHandlesWidePlaceholders ensures placeholders wider than the
// %03d minimum (i.e., 1000+ substitutions) are restored correctly rather than
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
