package ast

import (
	"reflect"
	"testing"
)

func TestModify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		il, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if il.Value != 1 {
			return node
		}

		il.Value = 2
		return il
	}

	tests := []struct {
		input    Node
		expected Node
	}{
		{one(), two()},
		{
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: one()},
				},
			},
			&Program{
				Statements: []Statement{
					&ExpressionStatement{Expression: two()},
				},
			},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: two()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&InfixExpression{Left: two(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "+", Right: one()},
			&PrefixExpression{Operator: "+", Right: two()},
		},
		{
			&IndexExpression{Left: one(), Index: one()},
			&IndexExpression{Left: two(), Index: two()},
		},
		{
			&IfExpression{
				Condition: one(),
				Then: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
				Else: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&IfExpression{
				Condition: two(),
				Then: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
				Else: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&ReturnStatement{ReturnValue: one()},
			&ReturnStatement{ReturnValue: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: one()},
					},
				},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body: &BlockStatement{
					Statements: []Statement{
						&ExpressionStatement{Expression: two()},
					},
				},
			},
		},
		{
			&ArrayLiteral{Elements: []Expression{one(), one()}},
			&ArrayLiteral{Elements: []Expression{two(), two()}},
		},
	}
	// Iterate over the test cases
	for _, tt := range tests {
		modified := Modify(tt.input, turnOneIntoTwo)

		equal := reflect.DeepEqual(modified, tt.expected)
		if !equal {
			t.Errorf("not equal. got=%#v, want=%#v", modified, tt.expected)
		}
	}

	// Hash
	hashLiteral := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			two(): two(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, value := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		if key.Value != 2 {
			t.Errorf("Value is not 2, got %d", key.Value)
		}

		val, _ := value.(*IntegerLiteral)
		if val.Value != 2 {
			t.Errorf("Value is not 2, got %d", val.Value)
		}
	}
}
