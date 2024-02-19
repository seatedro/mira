package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	ARRAY_TYPE    = "ARRAY"
	HASH_TYPE     = "HASH"
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

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_TYPE }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Bool) HashKey() HashKey {
	var hash uint64

	if b.Value {
		hash = 1
	} else {
		hash = 0
	}

	return HashKey{Type: b.Type(), Value: hash}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_TYPE }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}
