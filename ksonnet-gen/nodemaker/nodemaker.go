package nodemaker

import (
	"fmt"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/pkg/errors"
)

// Noder is an entity that can be converted to a jsonnet node.
type Noder interface {
	Node() ast.Node
}

type field struct {
	key   Key
	value Noder
}

// ObjectOptOneline is a functional option which sets the object's oneline status.
func ObjectOptOneline(oneline bool) ObjectOpt {
	return func(o *Object) {
		o.oneline = oneline
	}
}

// ObjectOpt is a functional option for Object.
type ObjectOpt func(*Object)

// Object is an item that can have multiple keys with values.
type Object struct {
	oneline bool
	fields  map[string]Noder
	keys    map[string]Key
	keyList []string
}

var _ Noder = (*Object)(nil)

// NewObject creates an Object. ObjectOpt functional arguments can be used to configure the
// newly generated key.
func NewObject(opts ...ObjectOpt) *Object {
	o := &Object{
		fields:  make(map[string]Noder),
		keys:    make(map[string]Key),
		keyList: make([]string, 0),
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}

// OnelineObject is a convenience method for creating a online object.
func OnelineObject(opts ...ObjectOpt) *Object {
	opts = append(opts, ObjectOptOneline(true))
	return NewObject(opts...)
}

// Set sets a field with a value.
func (o *Object) Set(key Key, value Noder) error {
	name := key.name

	if _, ok := o.keys[name]; ok {
		return errors.Errorf("field %q already exists in the object", name)
	}

	o.keys[name] = key
	o.fields[name] = value
	o.keyList = append(o.keyList, name)

	return nil
}

// Get retrieves a field by name.
func (o *Object) Get(keyName string) Noder {
	return o.fields[keyName]
}

// Node converts the object to a jsonnet node.
func (o *Object) Node() ast.Node {
	ao := &ast.Object{
		Oneline: o.oneline,
	}

	for _, name := range o.keyList {
		k := o.keys[name]
		v := o.fields[name]

		of := ast.ObjectField{
			Id:         newIdentifier(k.name),
			Kind:       k.category,
			Hide:       k.visibility,
			Expr2:      v.Node(),
			Comment:    o.generateComment(k.comment),
			Method:     k.Method(),
			SuperSugar: k.Mixin(),
		}

		ao.Fields = append(ao.Fields, of)
	}

	return ao
}

func (o *Object) generateComment(text string) *ast.Comment {
	if text != "" {
		return &ast.Comment{Text: text}
	}

	return nil
}

// StringDouble is double quoted string.
type StringDouble struct {
	text string
}

var _ Noder = (*StringDouble)(nil)

// NewStringDouble creates an instance of StringDouble.
func NewStringDouble(text string) *StringDouble {
	return &StringDouble{
		text: text,
	}
}

// Node converts the StringDouble to a jsonnet node.
func (t *StringDouble) Node() ast.Node {
	return &ast.LiteralString{
		Kind:  ast.StringDouble,
		Value: t.text,
	}
}

// Number is an a number.
type Number struct {
	number float64
}

var _ Noder = (*Number)(nil)

// NewNumber creates an instance of Number.
func NewNumber(number float64) *Number {
	return &Number{
		number: number,
	}
}

// Node converts the Number to a jsonnet node.
func (t *Number) Node() ast.Node {
	return &ast.LiteralNumber{
		Value: t.number,
	}
}

// Array is an an array.
type Array struct {
	elements []Noder
}

var _ Noder = (*Array)(nil)

// NewArray creates an instance of Array.
func NewArray(elements []Noder) *Array {
	return &Array{
		elements: elements,
	}
}

// Node converts the Array to a jsonnet node.
func (t *Array) Node() ast.Node {
	var nodes []ast.Node
	for _, element := range t.elements {
		nodes = append(nodes, element.Node())
	}

	return &ast.Array{
		Elements: nodes,
	}
}

// KeyOptCategory is a functional option for setting key category
func KeyOptCategory(kc ast.ObjectFieldKind) KeyOpt {
	return func(k *Key) {
		k.category = kc
	}
}

// KeyOptVisibility is a functional option for setting key visibility
func KeyOptVisibility(kv ast.ObjectFieldHide) KeyOpt {
	return func(k *Key) {
		k.visibility = kv
	}
}

// KeyOptComment is a functional option for setting a comment on a key
func KeyOptComment(text string) KeyOpt {
	return func(k *Key) {
		k.comment = text
	}
}

// KeyOptMixin is a functional option for setting this key as a mixin
func KeyOptMixin(b bool) KeyOpt {
	return func(k *Key) {
		k.mixin = b
	}
}

// KeyOptParams is functional option for setting params for a key. If there are no required
// parameters, pass an empty []string.
func KeyOptParams(params []string) KeyOpt {
	return func(k *Key) {
		k.params = params
	}
}

// KeyOpt is a functional option for configuring Key.
type KeyOpt func(k *Key)

// Key names a fields in an object.
type Key struct {
	name       string
	category   ast.ObjectFieldKind
	visibility ast.ObjectFieldHide
	comment    string
	params     []string
	mixin      bool
}

// NewKey creates an instance of Key. KeyOpt functional options can be used to configure the
// newly generated key.
func NewKey(name string, opts ...KeyOpt) Key {
	k := Key{
		name:     name,
		category: ast.ObjectFieldID,
	}
	for _, opt := range opts {
		opt(&k)
	}

	return k
}

// InheritedKey is a convenience method for creating an inherited key.
func InheritedKey(name string, opts ...KeyOpt) Key {
	opts = append(opts, KeyOptVisibility(ast.ObjectFieldInherit))
	return NewKey(name, opts...)
}

// LocalKey is a convenience method for creating a local key.
func LocalKey(name string, opts ...KeyOpt) Key {
	opts = append(opts, KeyOptCategory(ast.ObjectLocal))
	return NewKey(name, opts...)
}

// FunctionKey is a convenience method for creating a function key.
func FunctionKey(name string, args []string, opts ...KeyOpt) Key {
	opts = append(opts, KeyOptParams(args), KeyOptCategory(ast.ObjectFieldID))
	return NewKey(name, opts...)
}

// Method returns the jsonnet AST object file method parameter.
func (k Key) Method() *ast.Function {
	if k.params == nil {
		return nil
	}

	f := &ast.Function{
		Parameters: ast.Parameters{
			Required: ast.Identifiers{},
		},
	}

	for _, p := range k.params {
		f.Parameters.Required = append(f.Parameters.Required, *newIdentifier(p))
	}

	return f
}

// Mixin returns true if the jsonnet object should be super sugared.
func (k Key) Mixin() bool {
	return k.mixin
}

// BinaryOp is a binary operation.
type BinaryOp string

const (
	// BopPlus is +
	BopPlus BinaryOp = "+"
	// BopEqual is ==
	BopEqual = "=="
)

// Binary represents a binary operation
type Binary struct {
	Left  Noder
	Right Noder
	Op    BinaryOp
}

var _ Noder = (*Binary)(nil)

// NewBinary creates an instance of Binary.
func NewBinary(left, right Noder, op BinaryOp) *Binary {
	return &Binary{
		Left:  left,
		Right: right,
		Op:    op,
	}
}

// Node converts a BinaryOp into an ast node. This will panic if the binary operator
// is unknown.
func (b *Binary) Node() ast.Node {
	op, ok := ast.BopMap[string(b.Op)]
	if !ok {
		panic(fmt.Sprintf("%q is an invalid binary operation", b.Op))
	}

	return &ast.Binary{
		Left:  b.Left.Node(),
		Right: b.Right.Node(),
		Op:    op,
	}
}

// Var represents a variable.
type Var struct {
	ID string
}

var _ Noder = (*Binary)(nil)

// NewVar creates an instance of Var.
func NewVar(id string) *Var {
	return &Var{
		ID: id,
	}
}

// Node converts the var to a jsonnet ast node.
func (v *Var) Node() ast.Node {
	return &ast.Var{
		Id: *newIdentifier(v.ID),
	}
}

// Self represents self.
type Self struct{}

var _ Noder = (*Self)(nil)

// Node converts self to a jsonnet self node.
func (s *Self) Node() ast.Node {
	return &ast.Self{}
}

// Conditional represents a conditional
type Conditional struct {
	Cond        Noder
	BranchTrue  Noder
	BranchFalse Noder
}

var _ Noder = (*Conditional)(nil)

// NewConditional creates an instance of Conditional.
func NewConditional(cond, tbranch, fbranch Noder) *Conditional {
	return &Conditional{
		Cond:        cond,
		BranchTrue:  tbranch,
		BranchFalse: fbranch,
	}
}

// Node converts the Conditional to a jsonnet ast node.
func (c *Conditional) Node() ast.Node {
	cond := &ast.Conditional{
		Cond:       c.Cond.Node(),
		BranchTrue: c.BranchTrue.Node(),
	}

	if c.BranchFalse != nil {
		cond.BranchFalse = c.BranchFalse.Node()
	}

	return cond
}

// Apply represents an application of a function.
type Apply struct {
	target Noder
	args   []Noder
}

var _ Noder = (*Apply)(nil)

// NewApply creates an instance of Apply.
func NewApply(target Noder, args ...Noder) *Apply {
	return &Apply{
		target: target,
		args:   args,
	}
}

// ApplyCall creates an Apply using a method string.
func ApplyCall(method string, args ...Noder) *Apply {
	return NewApply(NewCall(method), args...)
}

// Node converts the Apply to a jsonnet ast node.
func (a *Apply) Node() ast.Node {
	nodes := make([]ast.Node, 0)
	for _, arg := range a.args {
		nodes = append(nodes, arg.Node())
	}

	apply := &ast.Apply{
		Target: a.target.Node(),
		Arguments: ast.Arguments{
			Positional: nodes,
		},
	}

	return apply
}

// Call is a function call.
type Call struct {
	parts []string
}

var _ Noder = (*Call)(nil)

// NewCall creates an instance of Call.
func NewCall(method string) *Call {
	parts := strings.Split(method, ".")

	return &Call{
		parts: parts,
	}
}

// Node converts the Call to a jsonnet ast node.
func (c *Call) Node() ast.Node {
	var head *ast.Index
	var cur *ast.Index

	if len(c.parts) == 1 {
		return NewVar(c.parts[0]).Node()
	}

	for i := 0; i < len(c.parts)-1; i++ {
		newIndex := &ast.Index{
			Id: newIdentifier(c.parts[i]),
		}
		if head == nil {
			head = newIndex
			cur = newIndex
		} else {
			cur.Target = newIndex
			cur = newIndex
		}
	}

	cur.Target = NewVar(c.parts[len(c.parts)-1]).Node()

	return head
}

// newIdentifier creates an identifier.
func newIdentifier(value string) *ast.Identifier {
	id := ast.Identifier(value)
	return &id
}
