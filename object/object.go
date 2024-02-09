package object

import (
	"bytes"
	"fmt"
	"mira/ast"
	"strings"
)

const (
	INTEGER_TYPE  = "INTEGER"
	BOOL_TYPE     = "BOOL"
	NULL_TYPE     = "NULL"
	RETURN_VALUE  = "RETURN_VALUE"
	ERROR_TYPE    = "ERROR"
	FUNCTION_TYPE = "FUNCTION"
	STRING_TYPE   = "STRING"
	BUILTIN_TYPE  = "BUILTIN"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_TYPE }

type Bool struct {
	Value bool
}

func (b *Bool) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Bool) Type() ObjectType { return BOOL_TYPE }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_TYPE }

type Null struct{}

func (i *Null) Inspect() string  { return "null" }
func (i *Null) Type() ObjectType { return NULL_TYPE }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_TYPE }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Body       *ast.BlockStatement
	Env        *Env
	Parameters []*ast.Identifier
}

func (f *Function) Type() ObjectType { return FUNCTION_TYPE }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type (
	BuiltinFunction func(args ...Object) Object
	Builtin         struct {
		Fn BuiltinFunction
	}
)

func (b *Builtin) Type() ObjectType { return BUILTIN_TYPE }
func (b *Builtin) Inspect() string  { return "builtin function" }
