// parser/parser_test.go

package parser

import (
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
