// Copyright (c) 2020, Alliance for Sustainable Energy, LLC.
// All rights reserved.

package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"strings"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/urbanopt/modelica-fmt/thirdparty/parser"
)

const (
	// indent
	spaceIndent = "  "
)

// insertIndentBefore returns true if the rule should be on a new line and indented
func (l *modelicaListener) insertIndentBefore(rule antlr.ParserRuleContext) bool {
	switch rule.(type) {
	case
		parser.IElementContext,
		parser.IEquationsContext,
		parser.IAlgorithm_statementsContext,
		parser.IControl_structure_bodyContext,
		parser.IAnnotationContext,
		parser.IExpression_listContext,
		parser.IConstraining_clauseContext,
		parser.IIf_expressionContext,
		parser.IIf_expression_bodyContext:
		return true
	case parser.IString_commentContext:
		return 0 == l.inAnnotation
	case
		parser.IArgumentContext,
		parser.INamed_argumentContext:
		return 0 == l.inAnnotation || 0 < l.inModelAnnotation
	case parser.IExpressionContext:
		if len(l.modelAnnotationVectorStack) == 0 {
			return false
		}

		// handle expression which is an element of a vector (array_arguments) and within model annotation
		arrayArgumentsNode, ok := rule.GetParent().(*parser.Array_argumentsContext)
		if !ok {
			return false
		}

		// check if the vector is the same as the one on top of our stack
		thisVectorInterval := arrayArgumentsNode.GetParent().(*parser.VectorContext).GetSourceInterval()
		stackVectorInterval := l.modelAnnotationVectorStack[len(l.modelAnnotationVectorStack)-1].GetSourceInterval()
		if thisVectorInterval.Start == stackVectorInterval.Start && thisVectorInterval.Stop == stackVectorInterval.Stop {
			return true
		}
		return false
	case parser.IFunction_argumentContext:
		return 0 == l.inNamedArgument && 0 == l.inVector && (0 == l.inAnnotation || 0 < l.inModelAnnotation)
	default:
		return false
	}
}

// insertSpaceBeforeToken returns true if a space should be inserted before the current token
func insertSpaceBeforeToken(currentTokenText, previousTokenText string) bool {
	switch currentTokenText {
	case "(":
		// add a space between 'annotation' and opening parens
		if previousTokenText == "annotation" {
			return true
		}
		fallthrough
	default:
		return !tokenInGroup(previousTokenText, noSpaceAfterTokens) &&
			!tokenInGroup(currentTokenText, noSpaceBeforeTokens)
	}
}

// insertNewlineBefore returns true if the rule should be on a new line
func insertNewlineBefore(rule antlr.ParserRuleContext) bool {
	switch rule.(type) {
	case
		parser.ICompositionContext,
		parser.IEquationsContext,
		parser.IIf_expression_conditionContext,
		parser.IElseif_expression_conditionContext,
		parser.IElse_expression_conditionContext:
		return true
	default:
		return false
	}
}

var (
	// tokens which should *generally* not have a space after them
	// this can be overridden in the insertSpace function
	noSpaceAfterTokens = []string{
		"(",
		"=",
		".",
		"[",
		"{",
		"-", "+", "^", "*", "/",
		";",
		",",
		":", // array range constructor
	}

	// tokens which should *generally* not have a space before them
	// this can be overridden in the insertSpace function
	noSpaceBeforeTokens = []string{
		"(", ")",
		"[", "]",
		"}",
		";",
		"=",
		",",
		".",
		"-", "+", "^", "*", "/",
		":", // array range constructor
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

type indent int

const (
	renderIndent indent = iota
	ignoreIndent
)

// modelicaListener is used to format the parse tree
type modelicaListener struct {
	*parser.BaseModelicaListener               // parser
	writer                       *bufio.Writer // writing destination
	indentationStack             []indent      // a stack used for tracking rendered and ignored indentations
	onNewLine                    bool          // true when write position succeeds a newline character
	lineIndentIncreased          bool          // true when the indentation level has already been increased for a line
	previousTokenText            string        // text of previous token
	previousTokenIdx             int           // index of previous token
	commentTokens                []antlr.Token // stores comments to insert while writing

	// modelAnnotationVectorStack is a stack which stores `vector` contexts,
	// which is used for conditionally indenting vector children
	// Specifically, a vector inside of a model annotation will have indented elements
	// if that vector has one or more elements which are function calls, class modifications or similar
	// (ie not if all elements are numbers, more vectors, etc)
	//
	// The last element of the slice is the first `vector` context ancestor whose contents
	// must be indented on new lines
	// For example, we would like model annotations to look like this:
	// annotation (
	// 	Abc(
	// 		paramA={
	// 			SomeIdentifier(
	// 				1,
	// 				2),
	//			123}))
	//
	// However, due to existing rules, we would end up with something like this
	// annotation (
	// 	Abc(
	// 		paramA={SomeIdentifier(
	// 			1,
	// 			2), 123}))
	//
	// Thus by pushing/popping vectors we can check if an expression in a vector
	// should be indented or not by checking if the top of the stack is its ancestor
	modelAnnotationVectorStack []antlr.RuleContext

	// NOTE: consider refactoring this simple approach for context awareness with
	// a set.
	// It should probably be map[string]int for rule name and current count (rules can be recursive, ie inside the same rule multiple times)
	inAnnotation      int // counts number of current or ancestor contexts that are annotation rule
	inModelAnnotation int // counts number of current or ancestor contexts that are model annotation rule
	inNamedArgument   int // counts number of current or ancestor contexts that are named argument
	inVector          int // counts number of current or ancestor contexts that are vector
}

func newListener(out io.Writer, commentTokens []antlr.Token) *modelicaListener {
	return &modelicaListener{
		BaseModelicaListener: &parser.BaseModelicaListener{},
		writer:               bufio.NewWriter(out),
		onNewLine:            true,
		lineIndentIncreased:  false,
		inAnnotation:         0,
		inModelAnnotation:    0,
		inVector:             0,
		inNamedArgument:      0,
		previousTokenText:    "",
		previousTokenIdx:     -1,
		commentTokens:        commentTokens,
	}
}

func (l *modelicaListener) close() {
	err := l.writer.Flush()
	if err != nil {
		panic(err)
	}
}

// indentation returns the writer's current number of *rendered* indentations
func (l *modelicaListener) indentation() int {
	nRenderIndents := 0
	for _, indentType := range l.indentationStack {
		if indentType == renderIndent {
			nRenderIndents++
		}
	}

	return nRenderIndents
}

// maybeIndent should be called when the writer's indentation is to be increased
func (l *modelicaListener) maybeIndent() {
	// Only increase indentation if it hasn't been changed already, otherwise ignore it
	// NOTE: This means that there can be at most 1 increase in indentation per line
	// This is a bit of a hack to avoid having an overindented line, occurring when
	// multiple rules want to be indented and we want it to be indented only once

	if !l.lineIndentIncreased {
		l.indentationStack = append(l.indentationStack, renderIndent)

		// WARNING: this is coupled with writeNewline, which should reset
		// lineIndentIncreased to false
		l.lineIndentIncreased = true
	} else {
		l.indentationStack = append(l.indentationStack, ignoreIndent)
	}
}

// maybeDedent should be called when the writer's indentation is to be decreased
func (l *modelicaListener) maybeDedent() {
	l.indentationStack = l.indentationStack[:len(l.indentationStack)-1]
}

func (l *modelicaListener) writeNewline() {
	l.writer.WriteString("\n")
	l.onNewLine = true

	// WARNING: this is coupled with maybeIndent, which uses this state
	l.lineIndentIncreased = false
}

func (l *modelicaListener) writeComment(comment antlr.Token) {
	l.writeSpaceBefore(comment)
	l.writer.WriteString(comment.GetText())
	if comment.GetTokenType() == parser.ModelicaLexerLINE_COMMENT {
		l.writeNewline()
	}
}

func (l *modelicaListener) writeSpaceBefore(token antlr.Token) {
	if l.onNewLine {
		// insert indentation
		if l.indentation() > 0 {
			indentation := l.indentation()
			l.writer.WriteString(strings.Repeat(spaceIndent, indentation))
		}
		l.onNewLine = false
	} else if insertSpaceBeforeToken(token.GetText(), l.previousTokenText) {
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

	if node.GetText() == ";" {
		l.writeNewline()
	}

	l.previousTokenText = node.GetText()
	l.previousTokenIdx = node.GetSymbol().GetTokenIndex()
}

func (l *modelicaListener) EnterEveryRule(node antlr.ParserRuleContext) {
	if insertNewlineBefore(node) && !l.onNewLine {
		l.writeNewline()
	}

	if l.insertIndentBefore(node) {
		if !l.onNewLine {
			l.writeNewline()
		}
		l.maybeIndent()
	}
}

func (l *modelicaListener) ExitEveryRule(node antlr.ParserRuleContext) {
	if l.insertIndentBefore(node) {
		l.maybeDedent()
	}
}

func (l *modelicaListener) EnterAnnotation(node *parser.AnnotationContext) {
	l.inAnnotation++
}

func (l *modelicaListener) ExitAnnotation(node *parser.AnnotationContext) {
	l.inAnnotation--
}

func (l *modelicaListener) EnterModel_annotation(node *parser.Model_annotationContext) {
	l.inModelAnnotation++
}

func (l *modelicaListener) ExitModel_annotation(node *parser.Model_annotationContext) {
	l.inModelAnnotation--
}

func (l *modelicaListener) EnterVector(node *parser.VectorContext) {
	l.inVector++
	if l.inModelAnnotation > 0 {
		// if this array uses an iterator for construction it gets no special treatment
		if _, ok := node.GetChild(0).(parser.Array_iterator_constructorContext); ok {
			return
		}

		// check if there is an element of this vector which would require indentation
		for _, child := range node.Array_arguments().GetChildren() {
			expressionNode, ok := child.(*parser.ExpressionContext)
			if !ok {
				continue
			}
			startToken := expressionNode.GetStart()
			if startToken.GetTokenType() == parser.ModelicaLexerIDENT {
				l.modelAnnotationVectorStack = append(l.modelAnnotationVectorStack, node)
				break
			}
		}
	}
}

func (l *modelicaListener) ExitVector(node *parser.VectorContext) {
	l.inVector--
	if len(l.modelAnnotationVectorStack) > 0 {
		annotationVectorInterval := l.modelAnnotationVectorStack[len(l.modelAnnotationVectorStack)-1].GetSourceInterval()
		thisVectorInterval := node.GetSourceInterval()
		if annotationVectorInterval.Start == thisVectorInterval.Start && annotationVectorInterval.Stop == thisVectorInterval.Stop {
			l.modelAnnotationVectorStack = l.modelAnnotationVectorStack[:len(l.modelAnnotationVectorStack)-1]
		}
	}
}

func (l *modelicaListener) EnterNamed_argument(node *parser.Named_argumentContext) {
	l.inNamedArgument++
}

func (l *modelicaListener) ExitNamed_argument(node *parser.Named_argumentContext) {
	l.inNamedArgument--
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
	if tokenType == parser.ModelicaLexerCOMMENT || tokenType == parser.ModelicaLexerLINE_COMMENT {
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
