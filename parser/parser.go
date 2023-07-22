// parser/parser.go

package parser

import (
	"fmt"
	"mira/ast"
	"mira/lexer"
	"mira/token"
// INFO: Operator Precedences
const (
	_ int = iota
	LOWEST
	EQUALS
	COMPARISON
	SUM
	PRODUCT
	PREFIX
	CALL
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
	errors    []string

	prefixParsers map[token.TokenType]prefixParseFn
	infixParsers  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParsers = make(map[token.TokenType]prefixParseFn)
	p.prefixParsers[token.IDENTIFIER] = p.parseIdentifier
	// Eg: let x = 5;
	// Calling twice because initially currToken = nil, nextToken = let.
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		stmnt := p.parseStatement()
		if stmnt != nil {
			program.Statements = append(program.Statements, stmnt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmnt := &ast.ExpressionStatement{Token: p.currToken}

	stmnt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParsers[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmnt := &ast.LetStatement{Token: p.currToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmnt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	//TODO: Skipping expressions for now

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmnt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	//TODO: Skipping expressions.

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}
