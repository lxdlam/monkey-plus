package object

import (
	"bytes"
	"fmt"
	"github.com/lxdlam/monkey-plus/ast"
	"hash/fnv"
	"strings"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
		return obj, ok
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Merge(other *Environment) {
	for k, v := range other.store {
		e.Set(k, v)
	}
}

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNC_OBJ         = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNC_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value     []uint8
	StringRep string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.StringRep }

func NewStringObject(raw string) *String {
	var value []uint8
	length := len(raw)

	for i := 0; i < length; i++ {
		if raw[i] == '\\' {
			var code uint8
			switch raw[i+1] {
			case 't':
				code = 9
			case 'b':
				code = 8
			case 'n':
				code = 10
			case 'r':
				code = 13
			case 'f':
				code = 12
			case '"':
				code = 34
			case '\\':
				code = 92
			}
			value = append(value, code)
			i++
		} else {
			value = append(value, raw[i])
		}
	}

	return &String{
		Value:     value,
		StringRep: raw,
	}
}

func (lhs *String) Concat(rhs *String) *String {
	leftLength := len(lhs.Value)
	newValue := make([]uint8, leftLength, leftLength)
	copy(newValue, lhs.Value)

	newValue = append(newValue, rhs.Value...)

	return &String{
		Value:     newValue,
		StringRep: string(newValue),
	}
}

func (lhs *String) Compare(rhs *String) int {
	leftLength := len(lhs.Value)
	rightLength := len(rhs.Value)

	for i := 0; i < leftLength && i < rightLength; i++ {
		if lhs.Value[i] != rhs.Value[i] {
			return int(lhs.Value[i]) - int(rhs.Value[i])
		}
	}

	if leftLength != rightLength {
		return leftLength - rightLength
	} else {
		return 0
	}
}

type BuiltinFunction func(env *Environment, args ...Object) Object

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

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
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
	Pairs  map[HashKey][]HashPair
	Length int
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pairList := range h.Pairs {
		if pairList != nil {
			for _, pair := range pairList {
				pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
			}
		}
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Helper functions for hash, because that the hash has complex operations
func (h *Hash) Len() int {
	return h.Length
}

func NewHash() *Hash {
	return &Hash{
		Pairs:  make(map[HashKey][]HashPair),
		Length: 0,
	}
}

// true if replace, otherwise create
func (h *Hash) Set(key, value Object) bool {
	hashKey := key.(Hashable).HashKey()
	switch key := key.(type) {
	case *String:
		pairs, ok := h.Pairs[hashKey]

		// If we don't found
		if !ok || pairs == nil {
			h.Pairs[hashKey] = []HashPair{{Key: key, Value: value}}
			h.Length++
			return false
		}

		// In order to resolve hash collision, we try to iterate the list
		length := len(pairs)
		var idx int
		for idx = 0; idx < length; idx++ {
			pairKey := pairs[idx].Key.(*String)

			if pairKey.Compare(key) == 0 {
				break
			}
		}

		// No collision, just append to the end
		if idx == length {
			h.Pairs[hashKey] = append(h.Pairs[hashKey], HashPair{Key: key, Value: value})
			h.Length++
			return false
		} else {
			h.Pairs[hashKey][idx].Value = value
			return true
		}
	default:
		pairs, ok := h.Pairs[hashKey]
		if !ok || pairs == nil {
			h.Pairs[hashKey] = []HashPair{{Key: key, Value: value}}
			h.Length++
			return false
		} else {
			h.Pairs[hashKey][0].Value = value
			return true
		}
	}
}

// true if successfully delete
func (h *Hash) Delete(key Object) bool {
	hashKey := key.(Hashable).HashKey()
	switch key := key.(type) {
	case *String:
		pairs, ok := h.Pairs[hashKey]

		if !ok || pairs == nil {
			return false
		}

		length := len(pairs)
		var idx int
		for idx = 0; idx < length; idx++ {
			pairKey := pairs[idx].Key.(*String)

			if pairKey.Compare(key) == 0 {
				newPairs := []HashPair{}
				newPairs = append(newPairs, pairs[:idx]...)
				newPairs = append(newPairs, pairs[idx+1:]...)
				h.Pairs[hashKey] = newPairs
				h.Length--
				return true
			}
		}

		return false
	default:
		pairs, ok := h.Pairs[hashKey]
		if !ok || pairs == nil {
			return false
		} else {
			delete(h.Pairs, hashKey)
			h.Length--
			return true
		}
	}
}

// Replace the real get
func (h *Hash) Get(key Object) (Object, bool) {
	hashKey := key.(Hashable)
	switch key := key.(type) {
	case *String:
		pairs, ok := h.Pairs[hashKey.HashKey()]
		if !ok || pairs == nil {
			return key, false
		}

		for _, pair := range pairs {
			pairKey := pair.Key.(*String)
			if pairKey.Compare(key) == 0 {
				return pair.Value, true
			}
		}

		return key, false
	default:
		pairs, ok := h.Pairs[hashKey.HashKey()]
		if !ok || pairs == nil {
			return key, false
		} else {
			return pairs[0].Value, true
		}
	}
}

func (h *Hash) Clone() *Hash {
	pairs := make(map[HashKey][]HashPair)
	for k, v := range h.Pairs {
		pairs[k] = v
	}

	return &Hash{
		Pairs:  pairs,
		Length: h.Length,
	}
}
