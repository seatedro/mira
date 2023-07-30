// parser/parser_test.go

package parser

import (
	"fmt"
	"mira/ast"
	"mira/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 12345;
`
	lex := lexer.New(input)
	parser := New(lex)

	program := parser.ParseProgram()
	checkParserErrors(t, parser)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("The number of statements in program is not 3. Got: %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmnt := program.Statements[i]
		if !testLetStatement(t, stmnt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not let. Got: %q", s.TokenLiteral())
		return false
	}

	letStmnt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s is not *ast.LetStatement. Got: %T", s)
		return false
	}

	if letStmnt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmnt.Name.Value)
		return false
	}

	if letStmnt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name, letStmnt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program doesn't have enough statements. Got: %d", len(program.Statements))
	}

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	ident, ok := stmnt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmnt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := New(lexer)
	program := parser.ParseProgram()
	checkParserErrors(t, parser)

	if len(program.Statements) != 1 {
		t.Fatalf("Program doesn't have enough statements. Got: %d", len(program.Statements))
	}

	stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	integer, ok := stmnt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmnt.Expression)
	}

	if integer.Value != 5 {
		t.Errorf("integer.Value not %d. got=%d", 5, integer.Value)
	}

	if integer.TokenLiteral() != "5" {
		t.Errorf("integer.TokenLiteral() not %s. got=%s", "5", integer.TokenLiteral())
	}
}

func TestPrefixOperatorExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!10;", "!", 10},
		{"-5;", "-", 5},
		{"--5;", "--", 5},
		{"++5;", "++", 5},
	}

	for _, tt := range prefixTests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program doesn't have enough statements. Got: %d", len(program.Statements))
		}

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		exp, ok := stmnt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmnt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Operator not %s, Got: %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestInfixOperatorExpression(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		leftValue  int64
		rightValue int64
	}{
		{"5 + 5;", "+", 5, 5},
		{"5 - 5;", "-", 5, 5},
		{"5 * 5;", "*", 5, 5},
		{"5 / 5;", "/", 5, 5},
		{"5 == 5;", "==", 5, 5},
		{"5 != 5;", "!=", 5, 5},
		{"5 > 5;", ">", 5, 5},
		{"5 < 5;", "<", 5, 5},
		{"5 >= 5;", ">=", 5, 5},
		{"5 <= 5;", "<=", 5, 5},
	}

	for _, tt := range prefixTests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program doesn't have enough statements. Got: %d", len(program.Statements))
		}

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		exp, ok := stmnt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmnt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Operator not %s, Got: %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func testInfixOperatorWithPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input          string
		prefixOperator string
		operator       string
		leftValue      int64
		rightValue     int64
	}{
		{"!5 != 5;", "!", "!=", 5, 5},
	}

	for _, tt := range prefixTests {
		lexer := lexer.New(tt.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParserErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program doesn't have enough statements. Got: %d", len(program.Statements))
		}

		stmnt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		exp, ok := stmnt.Expression.(*ast.InfixExpression)

		if !ok {
			t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmnt.Expression)
		}

		prefixVal := struct {
			operator     string
			integerValue int64
		}{
			operator:     tt.prefixOperator,
			integerValue: tt.leftValue,
		}

		if !testPrefixExp(t, exp, prefixVal) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Operator not %s, Got: %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}

}

func testPrefixExp(t *testing.T, exp *ast.InfixExpression, val struct {
	operator     string
	integerValue int64
}) bool {
	pe, ok := exp.Left.(*ast.PrefixExpression)

	if !ok {
		t.Fatalf("pe not *ast.PrefixExpression. got=%T", exp.Left)
	}

	if pe.Operator != val.operator {
		t.Fatalf("Operator not %s, Got: %s", val.operator, exp.Operator)
	}

	if !testIntegerLiteral(t, exp.Right, val.integerValue) {
		return false
	}

	return true
}
func testIntegerLiteral(t *testing.T, il ast.Expression, val int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il is not ast.IntegerLiteral, got: %T", il)
		return false
	}

	if integer.Value != val {
		t.Errorf("Value not equal to %d, Got: %d", val, integer.Value)
		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", val) {
		t.Errorf("integer.TokenLiteral not equal to %d, Got: %s", val, integer.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.errors
	fmt.Printf("")

	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
