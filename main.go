// Package main runs the formatter
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	parens = flag.Bool("p", false, "always break and indent contents in parens")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: modelicafmt [path ...]")
	flag.PrintDefaults()
}

func isModelicaFile(f os.FileInfo) bool {
	name := f.Name()
	return !f.IsDir() && !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".mo")
}

func visitFile(path string, f os.FileInfo, err error) error {
	if isModelicaFile(f) {
		err = processFile(path, os.Stdout)
	}

	if err != nil && !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return nil
}

func walkDir(path string) {
	filepath.Walk(path, visitFile)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "error: must provide at least one file or directory")
		os.Exit(2)
	}

	if *parens {
		alwaysIndentParens = true
	}

	for i := 0; i < flag.NArg(); i++ {
		path := flag.Arg(i)
		switch dir, err := os.Stat(path); {
		case err != nil:
			fmt.Fprintln(os.Stderr, "error: "+err.Error())
			os.Exit(2)
		case dir.IsDir():
			walkDir(path)
		default:
			if err := processFile(path, os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
}
