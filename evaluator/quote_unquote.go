package evaluator

import (
	"fmt"
	"mira/ast"
	"mira/object"
	"mira/token"
)

func quote(node ast.Node, env *object.Env) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Env) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertToASTNode(unquoted)
	})
}

func isUnquoteCall(node ast.Node) bool {
	if callExp, ok := node.(*ast.CallExpression); ok {
		return callExp.Function.TokenLiteral() == "unquote"
	}

	return false
}

func convertToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}
	case *object.Bool:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Bool{Token: t, Value: obj.Value}
	case *object.Quote:
		return obj.Node
	default:
		return nil
	}
}
