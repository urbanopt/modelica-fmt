package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/urbanopt/modelicafmt/thirdparty/parser"
)

var alwaysIndentParens = false

const (
	// lexer token types for comments
	commentTokenType     = 93
	lineCommentTokenType = 94

	// indent
	spaceIndent = "  "
)

// insertIndentBefore returns true if the rule should be indented
func insertIndentBefore(rule antlr.ParserRuleContext) bool {
	switch rule.(type) {
	case
			parser.IElementContext,
			parser.IEquationsContext,
			parser.IAlgorithm_statementsContext,
			parser.IControl_structure_bodyContext,
			parser.IString_commentContext,
			parser.IAnnotationContext,
			parser.IExpression_listContext:
		return true
	case
			parser.IArgumentContext,
			parser.INamed_argumentContext:
		return alwaysIndentParens
	default:
		return false
	}
}

// insertNewlineBefore returns true if the rule should be on a new line
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
		";",
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

// tokenInGroup returns true if a token is in a given list
func tokenInGroup(token string, group []string) bool {
	for _, other := range group {
		if token == other {
			return true
		}
	}
	return false
}

// modelicaListener is used to format the parse tree
type modelicaListener struct {
	*parser.BaseModelicaListener               // parser
	writer                       *bufio.Writer // writing destination
	indentation                  int           // current indentation depth (ignoring parens indentation)
	onNewLine                    bool          // true when write position succeeds a newline character
	previousTokenText            string        // text of previous token
	previousTokenIdx             int           // index of previous token
	numNestedParens              int           // tracks how deeply nested the write position is
	commentTokens                []antlr.Token // stores comments to insert while writing
}

func newListener(out io.Writer, commentTokens []antlr.Token) *modelicaListener {
	return &modelicaListener{
		BaseModelicaListener: &parser.BaseModelicaListener{},
		writer:               bufio.NewWriter(out),
		indentation:          0,
		onNewLine:            true,
		previousTokenText:    "",
		previousTokenIdx:     -1,
		numNestedParens:      0,
		commentTokens:        commentTokens,
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

func (l *modelicaListener) writeComment(comment antlr.Token) {
	l.writeSpaceBefore(comment)
	l.writer.WriteString(comment.GetText())
	if comment.GetTokenType() == lineCommentTokenType {
		l.writeNewline()
	}
}

func (l *modelicaListener) writeSpaceBefore(token antlr.Token) {
	if l.onNewLine {
		// insert indentation
		if l.indentation > 0 {
			indentation := l.indentation + l.numNestedParens
			l.writer.WriteString(strings.Repeat(spaceIndent, indentation))
		}
		l.onNewLine = false
	} else if !tokenInGroup(l.previousTokenText, noSpaceAfterTokens) && !tokenInGroup(token.GetText(), noSpaceBeforeTokens) {
		// insert a space
		l.writer.WriteString(" ")
	}
}

func (l *modelicaListener) VisitTerminal(node antlr.TerminalNode) {
	// if there's a comment that should go before this node, insert it first
	tokenIdx := node.GetSymbol().GetTokenIndex()
	for len(l.commentTokens) > 0 && tokenIdx > l.commentTokens[0].GetTokenIndex() && l.commentTokens[0].GetTokenIndex() > l.previousTokenIdx {
		commentToken := l.commentTokens[0]
		l.commentTokens = l.commentTokens[1:]
		l.writeComment(commentToken)
	}

	l.writeSpaceBefore(node.GetSymbol())

	l.writer.WriteString(node.GetText())

	if l.numNestedParens > 0 && node.GetText() == "," {
		l.writeNewline()
	}

	if node.GetText() == ";" {
		l.writeNewline()
	}

	l.previousTokenText = node.GetText()
	l.previousTokenIdx = node.GetSymbol().GetTokenIndex()
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

// commentCollector is a wrapper around the default lexer which collects comment
// tokens for later use
type commentCollector struct {
	antlr.TokenSource
	commentTokens []antlr.Token
}

func newCommentCollector(source antlr.TokenSource) commentCollector {
	return commentCollector{
		source,
		[]antlr.Token{},
	}
}

// NextToken returns the next token from the source
func (c *commentCollector) NextToken() antlr.Token {
	token := c.TokenSource.NextToken()

	tokenType := token.GetTokenType()
	if tokenType == commentTokenType || tokenType == lineCommentTokenType {
		c.commentTokens = append(c.commentTokens, token)
	}

	return token
}

// processFile formats a file
func processFile(filename string, out io.Writer) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	text := string(content)
	inputStream := antlr.NewInputStream(text)
	lexer := parser.NewModelicaLexer(inputStream)

	// quick runtime check for comment token types
	// TODO: figure out how to statically ensure this condition
	if lexer.SymbolicNames[commentTokenType] != "COMMENT" || lexer.SymbolicNames[lineCommentTokenType] != "LINE_COMMENT" {
		panic("Comment or line comment token types do not match")
	}

	// wrap the default lexer to collect comments and set it as the stream's source
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	tokenSource := newCommentCollector(lexer)
	stream.SetTokenSource(&tokenSource)

	p := parser.NewModelicaParser(stream)
	sd := p.Stored_definition()

	listener := newListener(out, tokenSource.commentTokens)
	defer listener.close()

	antlr.ParseTreeWalkerDefault.Walk(listener, sd)
	return nil
}
