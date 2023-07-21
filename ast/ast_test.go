// ast/ast_test.go
package ast

import (
	"mira/token"
	"testing"
)

func TestString(t *testing.T) {
	p := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x"},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "y"},
					Value: "y",
				},
			},
		},
	}

	if p.String() != "let x = y;" {
		t.Errorf("program.String() wrong. Got: %q", p.String())
	}
}
