package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/yurikdotdev/covfefescript/internal/ast"
)

type ObjectType string

const (
	MONEY_OBJ     = "MONEY"
	TRUTH_OBJ     = "TRUTH"
	COVFEFE_OBJ   = "COVFEFE"
	GIVE_ME_VALUE = "GIVE_ME"
	FUNCTION_OBJ  = "FUNCTION"
	TWEET_OBJ     = "TWEET"
	BUILTIN_OBJ   = "BUILTIN"
	ARRAY_OBJ     = "ARRAY"
	HASH_OBJ      = "HASH"
	ERROR_OBJ     = "ERROR"
)

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Money struct {
	Value int64
}

func (i *Money) Type() ObjectType { return MONEY_OBJ }
func (i *Money) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Money) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Truth struct {
	Value bool
}

func (b *Truth) Type() ObjectType { return TRUTH_OBJ }
func (b *Truth) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Truth) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

type Tweet struct {
	Value string
}

func (s *Tweet) Type() ObjectType { return TWEET_OBJ }
func (s *Tweet) Inspect() string  { return s.Value }
func (s *Tweet) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Covfefe struct{}

func (n *Covfefe) Type() ObjectType { return COVFEFE_OBJ }
func (n *Covfefe) Inspect() string  { return "COVFEFE" }

type BreakValue struct{}

func (bv *BreakValue) Type() ObjectType { return "BREAK_VALUE" }
func (bv *BreakValue) Inspect() string  { return "IT_WAS_RIGGED" }

type ContinueValue struct{}

func (cv *ContinueValue) Type() ObjectType { return "CONTINUE_VALUE" }
func (cv *ContinueValue) Inspect() string  { return "TIRED_OF_WINNING" }

var (
	YUGE     = &Truth{Value: true}
	LOSER    = &Truth{Value: false}
	COVFEFE  = &Covfefe{}
	BREAK    = &BreakValue{}
	CONTINUE = &ContinueValue{}
)

type GiveMeValue struct {
	Value Object
}

func (rv *GiveMeValue) Type() ObjectType { return GIVE_ME_VALUE }
func (rv *GiveMeValue) Inspect() string  { return rv.Value.Inspect() }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("MAKE_IT_BIG")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
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

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "SAD! " + e.Message }

func NewError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
