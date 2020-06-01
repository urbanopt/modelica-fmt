package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/urbanopt/modelicafmt/parser"
)

var alwaysIndentParens = false

func insertIndentBefore(rule antlr.ParserRuleContext) bool {
	switch rule.(type) {
	case
		parser.IElementContext,
		parser.IEquationsContext,
		parser.IAlgorithm_statementsContext,
		parser.IControl_structure_bodyContext,
		parser.IString_commentContext,
		parser.IAnnotationContext:
		return true
	default:
		return false
	}
}

func insertNewlineBefore(rule antlr.ParserRuleContext) bool {
	switch rule.(type) {
	case
		parser.ICompositionContext,
		parser.IEquationsContext:
		return true
	default:
		return false
	}
}

var (
	noSpaceAfterTokens = []string{
		"(",
		"=",
		".",
		"[",
		"{",
		"-", "^", "*", "/",
	}

	noSpaceBeforeTokens = []string{
		"(", ")",
		"[", "]",
		"}",
		";",
		"=",
		",",
		".",
		"-", "^", "*", "/",
	}

	listOpenTokens = []string{
		"(",
		"{",
	}

	listCloseTokens = []string{
		")",
		"}",
	}
)

func tokenInGroup(token string, group []string) bool {
	for _, other := range group {
		if token == other {
			return true
		}
	}
	return false
}

type modelicaListener struct {
	*parser.BaseModelicaListener
	writer            *bufio.Writer
	indentation       int
	onNewLine         bool
	previousTokenText string
	numNestedParens   int
}

func newListener(out io.Writer) *modelicaListener {
	return &modelicaListener{
		&parser.BaseModelicaListener{},
		bufio.NewWriter(out),
		0,
		true,
		"",
		0,
	}
}

func (l *modelicaListener) close() {
	err := l.writer.Flush()
	if err != nil {
		panic(err)
	}
}

func (l *modelicaListener) writeNewline() {
	l.writer.WriteString("\n")
	l.onNewLine = true
}

func (l *modelicaListener) VisitTerminal(node antlr.TerminalNode) {
	if alwaysIndentParens && tokenInGroup(node.GetText(), listCloseTokens) {
		l.numNestedParens--
	}

	if node.GetText() != ";" {
		if l.onNewLine {
			// insert indentation
			if l.indentation > 0 {
				indentation := l.indentation + l.numNestedParens
				l.writer.WriteString(strings.Repeat("    ", indentation))
			}
			l.onNewLine = false
		} else if !tokenInGroup(l.previousTokenText, noSpaceAfterTokens) && !tokenInGroup(node.GetText(), noSpaceBeforeTokens) {
			// insert a space
			l.writer.WriteString(" ")
		}
	}

	l.writer.WriteString(node.GetText())

	if alwaysIndentParens && tokenInGroup(node.GetText(), listOpenTokens) {
		l.writeNewline()
		l.numNestedParens++
	}

	if l.numNestedParens > 0 && node.GetText() == "," {
		l.writeNewline()
	}

	if node.GetText() == ";" {
		l.writeNewline()
	}

	l.previousTokenText = node.GetText()
}

func (l *modelicaListener) EnterEveryRule(node antlr.ParserRuleContext) {
	if insertNewlineBefore(node) {
		l.writeNewline()
	}
	if insertIndentBefore(node) {
		l.indentation++
		if !l.onNewLine {
			l.writeNewline()
		}
	}
}

func (l *modelicaListener) ExitEveryRule(node antlr.ParserRuleContext) {
	if insertIndentBefore(node) {
		l.indentation--
	}
}

func processFile(filename string, out io.Writer) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)
	inputStream := antlr.NewInputStream(text)
	lexer := parser.NewModelicaLexer(inputStream)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewModelicaParser(stream)

	listener := newListener(out)
	defer listener.close()
	antlr.ParseTreeWalkerDefault.Walk(listener, p.Stored_definition())

	return nil
}
