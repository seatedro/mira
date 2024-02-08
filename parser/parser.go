// parser/parser.go

package parser

import (
	"fmt"
	"mira/ast"
	"mira/lexer"
	"mira/token"
	"strconv"
)

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

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NEQ:      EQUALS,
	token.LT:       COMPARISON,
	token.GT:       COMPARISON,
	token.LE:       COMPARISON,
	token.GE:       COMPARISON,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l             *lexer.Lexer
	prefixParsers map[token.TokenType]prefixParseFn
	infixParsers  map[token.TokenType]infixParseFn
	currToken     token.Token
	peekToken     token.Token
	errors        []string
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

	// Infix Parse Functions
	p.infixParsers = make(map[token.TokenType]infixParseFn)
	p.infixParsers[token.EQ] = p.parseInfixExpression
	p.infixParsers[token.NEQ] = p.parseInfixExpression
	p.infixParsers[token.LT] = p.parseInfixExpression
	p.infixParsers[token.GT] = p.parseInfixExpression
	p.infixParsers[token.LE] = p.parseInfixExpression
	p.infixParsers[token.GE] = p.parseInfixExpression
	p.infixParsers[token.PLUS] = p.parseInfixExpression
	p.infixParsers[token.MINUS] = p.parseInfixExpression
	p.infixParsers[token.ASTERISK] = p.parseInfixExpression
	p.infixParsers[token.SLASH] = p.parseInfixExpression
	p.infixParsers[token.LPAREN] = p.parseCallExpression

	// Prefix Parse Functions
	p.prefixParsers = make(map[token.TokenType]prefixParseFn)
	p.prefixParsers[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParsers[token.INT] = p.parseIntegerLiteral
	p.prefixParsers[token.BANG] = p.parsePrefixExpression
	p.prefixParsers[token.MINUS] = p.parsePrefixExpression
	p.prefixParsers[token.DEC] = p.parsePrefixExpression
	p.prefixParsers[token.INC] = p.parsePrefixExpression
	p.prefixParsers[token.TRUE] = p.parseBoolean
	p.prefixParsers[token.FALSE] = p.parseBoolean
	p.prefixParsers[token.LPAREN] = p.parseGroupedExpression
	p.prefixParsers[token.IF] = p.parseIfExpression
	p.prefixParsers[token.FUNCTION] = p.parseFunctionLiteral
	p.prefixParsers[token.STRING] = p.parseStringLiteral

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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Bool{Token: p.currToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Could not parse %s as integer", p.currToken.Literal)
		p.errors = append(p.errors, errMsg)

		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fn := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFunctionParams()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{Token: p.currToken, Operator: p.currToken.Literal, Left: left}

	precedence := p.currPrecedence()
	p.nextToken()

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.currToken}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Then = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Else = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseFunctionParams() []*ast.Identifier {
	idents := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return idents
	}

	p.nextToken()

	// We do this here in case there is only one param
	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	idents = append(idents, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		// currToken = next param
		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		idents = append(idents, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return idents
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmnt := p.parseStatement()
		if stmnt != nil {
			block.Statements = append(block.Statements, stmnt)
		}
		p.nextToken()
	}

	return block
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
		msg := fmt.Sprintf("No prefix parse fn found for %s", p.currToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

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

	p.nextToken()

	stmnt.Value = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmnt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmnt := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	stmnt.ReturnValue = p.parseExpression(LOWEST)

	for p.peekTokenIs(token.SEMICOLON) {
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
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
