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

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
