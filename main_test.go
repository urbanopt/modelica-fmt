package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsModelicaFileIncludesTemplateExtensions(t *testing.T) {
	a := require.New(t)
	dir := t.TempDir()

	for _, name := range []string{"Model.mo", "Template.mot", "Template.mopt"} {
		filePath := path.Join(dir, name)
		a.NoError(os.WriteFile(filePath, []byte(""), 0644))
		info, err := os.Stat(filePath)
		a.NoError(err)
		a.True(isModelicaFile(info), "%s should be formatted", name)
	}
}
