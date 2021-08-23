package main

import (
	"bytes"
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
	{"gmt-coolingtower.mo", "gmt-coolingtower-out.mo", Config{-1, false}},
	{"functions.mo", "functions-out.mo", Config{-1, false}},
	{"example-no-within.mo", "example-no-within-out.mo", Config{-1, false}},
	{"example-arrays.mo", "example-arrays-out.mo", Config{-1, false}},
	{"gmt-building.mo", "gmt-building-out.mo", Config{-1, false}},
	{"gmt-building.mo", "gmt-building-80-out.mo", Config{80, false}},
	{"gmt-building.mo", "gmt-building-empty-lines-out.mo", Config{-1, true}},
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
