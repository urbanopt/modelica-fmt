package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const outputDir = "test_output"

func TestMain(m *testing.M) {
	_ = os.Mkdir(outputDir, 0755)
	code := m.Run()
	os.Exit(code)
}

func diffFiles(a, b string) (string, error) {
	cmd := exec.Command("git", "diff", "--no-index", a, b)
	var out bytes.Buffer
	cmd.Stdout = &out
	_ = cmd.Run()
	return out.String(), nil
}

var exampleFileTests = []struct {
	sourceFile      string
	outFile         string
	formatterConfig Config
}{
	{"gmt-coolingtower.mo", "gmt-coolingtower-out.mo", Config{-1, false, false}},
	{"functions.mo", "functions-out.mo", Config{-1, false, false}},
	{"example-no-within.mo", "example-no-within-out.mo", Config{-1, false, false}},
	{"example-arrays.mo", "example-arrays-out.mo", Config{-1, false, false}},
	{"gmt-building.mo", "gmt-building-out.mo", Config{-1, false, false}},
	{"gmt-building.mo", "gmt-building-80-out.mo", Config{80, false, false}},
	{"gmt-building.mo", "gmt-building-empty-lines-out.mo", Config{-1, true, false}},
	{"html-annotation.mo", "html-annotation-out.mo", Config{-1, false, false}},
	{"html-annotation.mo", "html-annotation-html-out.mo", Config{-1, false, true}},
	{"gmt-building.mo", "gmt-building-html-out.mo", Config{-1, false, true}},
}

func TestFormattingExamples(t *testing.T) {
	a := require.New(t)
	for _, testCase := range exampleFileTests {
		t.Run(testCase.sourceFile, func(t *testing.T) {
			// Setup
			testSourceFile := path.Join("examples", testCase.sourceFile)
			expectedOutFile := path.Join("examples", testCase.outFile)
			actualOutFile := path.Join(outputDir, testCase.outFile)
			file, err := os.Create(actualOutFile)
			a.NoError(err)
			defer file.Close()

			// Act
			err = processFile(
				testSourceFile,
				file,
				testCase.formatterConfig,
			)

			// Assert
			a.NoError(err)
			diff, err := diffFiles(expectedOutFile, actualOutFile)
			a.NoError(err)
			a.Len(diff, 0, "File diff should be empty")
		})
	}
}

// TestProcessFileReturnsErrorOnInvalidInput ensures that processFile returns an
// error when the input cannot be parsed cleanly (e.g. it contains an
// unrecognized token). This allows callers to avoid overwriting the original
// file with malformed output (see issue #34).
func TestProcessFileReturnsErrorOnInvalidInput(t *testing.T) {
	a := require.New(t)

	invalidContent := "model Test\n  Real x = 1 # bad token;\nequation\n  x = 2;\nend Test;\n"
	sourceFile := path.Join(outputDir, "invalid-input.mo")
	a.NoError(ioutil.WriteFile(sourceFile, []byte(invalidContent), 0644))
	defer os.Remove(sourceFile)

	var out bytes.Buffer
	err := processFile(sourceFile, &out, Config{-1, false, false})

	a.Error(err, "processFile should return an error for input with unrecognized tokens")
}
