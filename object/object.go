package object

import "fmt"

const (
	INTEGER_TYPE = "INTEGER"
	BOOL_TYPE    = "BOOL"
	NULL_TYPE    = "NULL"
	RETURN_VALUE = "RETURN_VALUE"
	ERROR_TYPE   = "ERROR"
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

func (i *Bool) Inspect() string  { return fmt.Sprintf("%t", i.Value) }
func (i *Bool) Type() ObjectType { return BOOL_TYPE }

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
