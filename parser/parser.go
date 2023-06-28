// parser/parser.go

package parser

import (
	"mira/ast"
	"mira/lexer"
	"mira/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Eg: let x = 5;
	// Calling twice because initially currToken = nil, nextToken = let.
	p.nextToken()
	p.nextToken()
	return p
}

func (p Parser) ParseProgram() *ast.Program {
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		break
	}
	return nil
}

func (p Parser) parseLetStatement() ast.Statement {
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

func (p Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

func (p Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return true
}
