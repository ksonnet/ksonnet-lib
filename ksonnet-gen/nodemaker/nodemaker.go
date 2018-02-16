package nodemaker

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/astext"
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
		o.Oneline = oneline
	}
}

// ObjectOpt is a functional option for Object.
type ObjectOpt func(*Object)

// Object is an item that can have multiple keys with values.
type Object struct {
	Oneline bool
	fields  map[string]Noder
	keys    map[string]Key
	keyList []string
}

var _ Noder = (*Object)(nil)

// KVFromMap creates a shallow object using a map.
func KVFromMap(m map[string]interface{}) (*Object, error) {
	if m == nil {
		return nil, errors.New("map is nil")
	}

	var names []string
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)

	o := NewObject()

	for _, name := range names {
		switch t := m[name].(type) {
		case string:
			o.Set(InheritedKey(name), NewStringDouble(t))
		case float64:
			o.Set(InheritedKey(name), NewFloat(t))
		case int:
			o.Set(InheritedKey(name), NewInt(t))
		case bool:
			o.Set(InheritedKey(name), NewBoolean(t))
		default:
			return nil, errors.Errorf("unsupported type %T", t)
		}
	}

	return o, nil
}

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
	ao := &astext.Object{
		Oneline: o.Oneline,
	}

	for _, name := range o.keyList {
		k := o.keys[name]
		v := o.fields[name]

		of := astext.ObjectField{
			Comment: o.generateComment(k.comment),
		}
		of.Id = newIdentifier(k.name)
		of.Kind = k.category
		of.Hide = k.visibility
		of.Expr2 = v.Node()
		of.Method = k.Method()
		of.SuperSugar = k.Mixin()

		ao.Fields = append(ao.Fields, of)
	}

	return ao
}

func (o *Object) generateComment(text string) *astext.Comment {
	if text != "" {
		return &astext.Comment{Text: text}
	}

	return nil
}

// Boolean is a boolean.
type Boolean struct {
	value bool
}

// NewBoolean creates an instance of Boolean.
func NewBoolean(value bool) *Boolean {
	return &Boolean{
		value: value,
	}
}

// Node converts Boolean to a jsonnet node.
func (b *Boolean) Node() ast.Node {
	return &ast.LiteralBoolean{
		Value: b.value,
	}
}

// StringDouble is double quoted string.
type StringDouble struct {
	text string
}

// NewStringDouble creates an instance of StringDouble.
func NewStringDouble(text string) *StringDouble {
	return &StringDouble{
		text: text,
	}
}

func (t *StringDouble) node() *ast.LiteralString {
	return &ast.LiteralString{
		Kind:  ast.StringDouble,
		Value: t.text,
	}
}

// Node converts the StringDouble to a jsonnet node.
func (t *StringDouble) Node() ast.Node {
	return t.node()
}

// Number is an a number.
type Number struct {
	number float64
	value  string
}

var _ Noder = (*Number)(nil)

// NewInt creates an integer number.
func NewInt(i int) *Number {
	return &Number{
		number: float64(i),
		value:  strconv.Itoa(i),
	}
}

// NewFloat creates a float instance of a number.
func NewFloat(f float64) *Number {
	return &Number{
		number: f,
		value:  strconv.FormatFloat(f, 'f', -1, 64),
	}
}

// Node converts the Number to a jsonnet node.
func (t *Number) Node() ast.Node {
	return &ast.LiteralNumber{
		Value:          t.number,
		OriginalString: t.value,
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

// KeyOptNamedParams is a functional option for setting named params for a key.
func KeyOptNamedParams(params ...OptionalArg) KeyOpt {
	return func(k *Key) {
		k.namedParams = params
	}
}

// KeyOpt is a functional option for configuring Key.
type KeyOpt func(k *Key)

var (
	jsonnetReservedWords = []string{"assert", "else", "error", "false", "for", "function", "if",
		"import", "importstr", "in", "local", "null", "tailstrict", "then", "self", "super", "true"}
)

// Key names a fields in an object.
type Key struct {
	name        string
	category    ast.ObjectFieldKind
	visibility  ast.ObjectFieldHide
	comment     string
	params      []string
	namedParams []OptionalArg
	mixin       bool
}

// NewKey creates an instance of Key. KeyOpt functional options can be used to configure the
// newly generated key.
func NewKey(name string, opts ...KeyOpt) Key {

	category := ast.ObjectFieldID
	for _, s := range jsonnetReservedWords {
		if s == name {
			category = ast.ObjectFieldStr
		}
	}

	k := Key{
		name:     name,
		category: category,
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
func (k *Key) Method() *ast.Function {
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

	for _, p := range k.namedParams {
		f.Parameters.Optional = append(f.Parameters.Optional, p.NamedParameter())
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
	// BopGreater is >
	BopGreater = ">"
	// BopAnd is &&
	BopAnd = "&&"
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

// OptionalArg is an optional argument.
type OptionalArg struct {
	Name    string
	Default Noder
}

// NamedArgument converts the OptionalArgument to a jsonnet NamedArgument.
func (oa *OptionalArg) NamedArgument() ast.NamedArgument {
	na := ast.NamedArgument{
		Name: *newIdentifier(oa.Name),
	}

	if oa.Default == nil {
		na.Arg = NewStringDouble("").Node()
	} else {
		na.Arg = oa.Default.Node()
	}

	return na
}

// NamedParameter converts the OptionalArgument to a jsonnet NamedParameter.
func (oa *OptionalArg) NamedParameter() ast.NamedParameter {
	np := ast.NamedParameter{
		Name: *newIdentifier(oa.Name),
	}

	if oa.Default != nil {
		np.DefaultArg = oa.Default.Node()
	}

	return np
}

// Apply represents an application of a function.
type Apply struct {
	target         Noder
	positionalArgs []Noder
	optionalArgs   []OptionalArg
}

var _ Noder = (*Apply)(nil)

// NewApply creates an instance of Apply.
func NewApply(target Noder, positionalArgs []Noder, optionalArgs []OptionalArg) *Apply {
	return &Apply{
		target:         target,
		positionalArgs: positionalArgs,
		optionalArgs:   optionalArgs,
	}
}

// ApplyCall creates an Apply using a method string.
func ApplyCall(method string, args ...Noder) *Apply {
	return NewApply(NewCall(method), args, nil)
}

// Node converts the Apply to a jsonnet ast node.
func (a *Apply) Node() ast.Node {
	apply := &ast.Apply{
		Target: a.target.Node(),
	}

	for _, arg := range a.positionalArgs {
		apply.Arguments.Positional = append(apply.Arguments.Positional, arg.Node())
	}

	for _, arg := range a.optionalArgs {
		apply.Arguments.Named = append(apply.Arguments.Named, arg.NamedArgument())
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

	parts := c.parts

	if len(parts) == 1 {
		return NewVar(parts[0]).Node()
	}

	for i := range parts {
		i = len(parts) - 1 - i
		if i == 0 {
			v := NewVar(parts[i]).Node()
			if head == nil {
				return v
			}

			cur.Target = v
			return head
		}

		newIndex := &ast.Index{
			Id: newIdentifier(parts[i]),
		}
		if head == nil {
			head = newIndex
			cur = newIndex
		} else {
			cur.Target = newIndex
			cur = newIndex
		}
	}

	return head
}

// Local is a local declaration.
type Local struct {
	name  string
	value Noder
	body  Noder
}

var _ Noder = (*Local)(nil)

// NewLocal creates an instance of Local.
func NewLocal(name string, value, body Noder) *Local {
	return &Local{name: name, value: value, body: body}
}

// Node converts the Local to a jsonnet ast node.
func (l *Local) Node() ast.Node {
	id := *newIdentifier(l.name)

	local := &ast.Local{
		Binds: ast.LocalBinds{
			{
				Variable: id,
				Body:     l.value.Node(),
			},
		},
	}

	if l.body != nil {
		local.Body = l.body.Node()
	}

	return local
}

// Import is an import declaration.
type Import struct {
	name string
}

var _ Noder = (*Import)(nil)

// NewImport creates an instance of Import.
func NewImport(name string) *Import {
	return &Import{name: name}
}

// Node converts the Import to a jsonnet ast node.
func (i *Import) Node() ast.Node {
	file := NewStringDouble(i.name)

	return &ast.Import{
		File: file.node(),
	}
}

// Function is a function.
type Function struct {
	req  []string
	body Noder
}

var _ Noder = (*Function)(nil)

// NewFunction creates an instance of Function.
func NewFunction(req []string, body Noder) *Function {
	return &Function{
		req:  req,
		body: body,
	}
}

// Node converts the Function to a jsonnet ast node.
func (f *Function) Node() ast.Node {
	fun := &ast.Function{
		Parameters: ast.Parameters{},
		Body:       f.body.Node(),
	}

	var ids ast.Identifiers
	for _, param := range f.req {
		ids = append(ids, *newIdentifier(param))
	}
	fun.Parameters.Required = ids

	return fun
}

// Combine combines multiple nodes into a single node. If one argument is passed,
// it is returned. If two or more arguments are passed, they are combined using a
// Binary.
func Combine(nodes ...Noder) Noder {
	l := len(nodes)

	switch {
	case l == 1:
		return nodes[0]
	case l >= 2:
		sum := NewBinary(nodes[0], nodes[1], BopPlus)

		for i := 2; i < l; i++ {
			sum = NewBinary(sum, nodes[i], BopPlus)
		}

		return sum
	}

	return NewObject()
}

// newIdentifier creates an identifier.
func newIdentifier(value string) *ast.Identifier {
	id := ast.Identifier(value)
	return &id
}
