package nodemaker

import (
	"fmt"
	"os"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/stretchr/testify/require"
)

func ExampleApply() {
	o := NewObject()
	k := LocalKey("foo")

	arg1 := NewStringDouble("arg1")

	a := ApplyCall("alpha.beta.charlie", arg1)

	if err := o.Set(k, a); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	if err := printer.Fprint(os.Stdout, o.Node()); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	// Output:
	// {
	//   local foo = alpha.beta.charlie("arg1"),
	// }
}

func ExampleArray() {
	o := NewObject()

	nodes := []Noder{NewStringDouble("hello")}

	t := NewArray(nodes)
	k := InheritedKey("foo")
	if err := o.Set(k, t); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	if err := printer.Fprint(os.Stdout, o.Node()); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	// Output:
	// {
	//   foo: ["hello"],
	// }
}

func ExampleBinary() {
	o := NewObject()
	k := NewKey("foo")
	b := NewBinary(NewVar("alpha"), NewVar("beta"), BopPlus)

	if err := o.Set(k, b); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	if err := printer.Fprint(os.Stdout, o.Node()); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	// Output:
	// {
	//   foo:: alpha + beta,
	// }
}

func ExampleObject() {
	o := NewObject()

	k := NewKey("foo")
	o2 := NewObject()

	if err := o.Set(k, o2); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	if err := printer.Fprint(os.Stdout, o.Node()); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	// Output:
	// {
	//   foo:: {
	//   },
	// }
}

func ExampleConditional() {
	o := NewObject()
	k := NewKey("foo")

	cond := NewBinary(
		NewVar("alpha"),
		NewVar("beta"),
		BopEqual,
	)
	trueBranch := NewObject()
	trueBranch.Set(
		NewKey(
			"foo",
			KeyOptCategory(ast.ObjectFieldID),
			KeyOptVisibility(ast.ObjectFieldInherit)),
		NewStringDouble("1"),
	)

	falseBranch := NewObject()
	falseBranch.Set(
		NewKey(
			"foo",
			KeyOptCategory(ast.ObjectFieldID),
			KeyOptVisibility(ast.ObjectFieldInherit)),
		NewStringDouble("2"),
	)

	c := NewConditional(cond, trueBranch, falseBranch)

	if err := o.Set(k, c); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	if err := printer.Fprint(os.Stdout, o.Node()); err != nil {
		fmt.Printf("error: %#v\n", err)
	}

	// Output:
	// {
	//   foo:: if alpha == beta then {
	//     foo: "1",
	//   } else {
	//     foo: "2",
	//   },
	// }
}

func TestObject(t *testing.T) {
	cases := []struct {
		name     string
		object   *Object
		expected ast.Node
	}{
		{
			name:     "empty",
			object:   NewObject(),
			expected: &ast.Object{},
		},
		{
			name:     "with a single key",
			object:   objectWithKeys(),
			expected: expectations["objectWithKeys"],
		},
		{
			name:     "inline",
			object:   inline(),
			expected: expectations["inline"],
		},
		{
			name:     "local field",
			object:   localField(),
			expected: expectations["localField"],
		},
		{
			name:     "text field",
			object:   textField(),
			expected: expectations["textField"],
		},
		{
			name:     "mixin field",
			object:   mixinField(),
			expected: expectations["mixinField"],
		},
		{
			name:     "number field",
			object:   numberField(),
			expected: expectations["numberField"],
		},
		{
			name:     "self field",
			object:   selfField(),
			expected: expectations["selfField"],
		},
		{
			name:     "array field",
			object:   arrayField(),
			expected: expectations["arrayField"],
		},
		{
			name:     "comment field",
			object:   commentedField(),
			expected: expectations["commentedField"],
		},
		{
			name:     "function",
			object:   functionField(),
			expected: expectations["function"],
		},
		{
			name:     "function with args",
			object:   functionFieldArg(),
			expected: expectations["functionWithArgs"],
		},
		{
			name:     "binary operation",
			object:   binaryOp(),
			expected: expectations["binaryOp"],
		},
		{
			name:     "conditional",
			object:   conditional(),
			expected: expectations["conditional"],
		},
		{
			name:     "conditional without false branch",
			object:   conditionalNoFalse(),
			expected: expectations["conditionalNoFalse"],
		},
		{
			name:     "local apply",
			object:   localApply(),
			expected: expectations["localApply"],
		},
		{
			name:     "local apply with an index",
			object:   localApply2(),
			expected: expectations["localApply2"],
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := tc.object.Node()

			require.Equal(t, tc.expected, node)
			// if !reflect.DeepEqual(tc.expected, node) {
			// 	t.Errorf("object didn't convert to expected node")
			// }
		})
	}

}

func objectWithKeys() *Object {
	o := NewObject()

	k := NewKey("foo")
	o2 := NewObject()

	o.Set(k, o2)

	return o
}

func inline() *Object {
	o := OnelineObject()

	k := NewKey("foo")
	o2 := NewObject()

	o.Set(k, o2)

	return o
}

func localField() *Object {
	o := NewObject()
	k := LocalKey("foo")
	o.Set(k, NewObject())
	return o
}

func sortedFields() *Object {
	o := NewObject()
	o.Set(NewKey("zField"), NewObject())
	o.Set(NewKey("aField"), NewObject())
	o.Set(LocalKey("aLocal"), NewObject())
	return o
}

func textField() *Object {
	o := NewObject()
	t := NewStringDouble("bar")

	k := InheritedKey("foo")
	o.Set(k, t)
	return o
}

func mixinField() *Object {
	o := NewObject()
	t := NewStringDouble("bar")
	k := NewKey("foo",
		KeyOptVisibility(ast.ObjectFieldInherit),
		KeyOptMixin(true))
	o.Set(k, t)
	return o
}

func numberField() *Object {
	o := NewObject()
	t := NewNumber(1)
	k := InheritedKey("foo")
	o.Set(k, t)
	return o
}

func selfField() *Object {
	o := NewObject()
	t := &Self{}
	k := InheritedKey("foo")
	o.Set(k, t)
	return o
}

func arrayField() *Object {
	o := NewObject()

	nodes := []Noder{NewStringDouble("hello")}

	t := NewArray(nodes)
	k := InheritedKey("foo")
	o.Set(k, t)
	return o
}

func commentedField() *Object {
	o := NewObject()

	k := NewKey("foo", KeyOptComment("a comment"))
	o2 := NewObject()

	o.Set(k, o2)

	return o
}

func functionField() *Object {
	o := NewObject()
	k := FunctionKey("foo", []string{})

	o.Set(k, NewObject())

	return o
}

func functionFieldArg() *Object {
	o := NewObject()
	k := FunctionKey("foo", []string{"arg1"})

	o.Set(k, NewObject())

	return o
}

func binaryOp() *Object {
	o := NewObject()
	k := NewKey("foo")
	b := NewBinary(NewVar("alpha"), NewVar("beta"), BopPlus)

	o.Set(k, b)

	return o
}

func conditional() *Object {
	o := NewObject()
	k := NewKey("foo")

	cond := NewBinary(
		NewVar("alpha"),
		NewVar("beta"),
		BopEqual,
	)
	trueBranch := NewObject()
	trueBranch.Set(
		NewKey(
			"foo",
			KeyOptCategory(ast.ObjectFieldID),
			KeyOptVisibility(ast.ObjectFieldInherit)),
		NewStringDouble("1"),
	)

	falseBranch := NewObject()
	falseBranch.Set(
		NewKey(
			"foo",
			KeyOptCategory(ast.ObjectFieldID),
			KeyOptVisibility(ast.ObjectFieldInherit)),
		NewStringDouble("2"),
	)

	c := NewConditional(cond, trueBranch, falseBranch)

	o.Set(k, c)

	return o
}

func conditionalNoFalse() *Object {
	o := NewObject()
	k := NewKey("foo")

	cond := NewBinary(
		NewVar("alpha"),
		NewVar("beta"),
		BopEqual,
	)
	trueBranch := NewObject()
	trueBranch.Set(
		NewKey(
			"foo",
			KeyOptCategory(ast.ObjectFieldID),
			KeyOptVisibility(ast.ObjectFieldInherit)),
		NewStringDouble("1"),
	)

	c := NewConditional(cond, trueBranch, nil)

	o.Set(k, c)

	return o
}

func localApply() *Object {
	o := NewObject()
	k := LocalKey("foo")

	call := NewCall("alpha")
	arg1 := NewStringDouble("arg1")
	a := NewApply(call, arg1)

	o.Set(k, a)
	return o
}

func localApply2() *Object {
	o := NewObject()
	k := LocalKey("foo")

	arg1 := NewStringDouble("arg1")
	a := ApplyCall("alpha.beta.charlie", arg1)

	o.Set(k, a)
	return o
}

var (
	expectations = map[string]ast.Node{
		"objectWithKeys": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:    newIdentifier("foo"),
					Hide:  ast.ObjectFieldHidden,
					Kind:  ast.ObjectFieldID,
					Expr2: &ast.Object{},
				},
			},
		},
		"inline": &ast.Object{
			Oneline: true,
			Fields: ast.ObjectFields{
				{
					Id:    newIdentifier("foo"),
					Hide:  ast.ObjectFieldHidden,
					Kind:  ast.ObjectFieldID,
					Expr2: &ast.Object{},
				},
			},
		},
		"localField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:    newIdentifier("foo"),
					Kind:  ast.ObjectLocal,
					Hide:  ast.ObjectFieldHidden,
					Expr2: &ast.Object{},
				},
			},
		},
		"textField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldInherit,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.LiteralString{
						Kind:  ast.StringDouble,
						Value: "bar",
					},
				},
			},
		},
		"mixinField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					SuperSugar: true,
					Id:         newIdentifier("foo"),
					Hide:       ast.ObjectFieldInherit,
					Kind:       ast.ObjectFieldID,
					Expr2: &ast.LiteralString{
						Kind:  ast.StringDouble,
						Value: "bar",
					},
				},
			},
		},
		"numberField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldInherit,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.LiteralNumber{
						Value: 1,
					},
				},
			},
		},
		"selfField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:    newIdentifier("foo"),
					Hide:  ast.ObjectFieldInherit,
					Kind:  ast.ObjectFieldID,
					Expr2: &ast.Self{},
				},
			},
		},
		"arrayField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldInherit,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.Array{
						Elements: []ast.Node{
							&ast.LiteralString{
								Kind:  ast.StringDouble,
								Value: "hello",
							},
						},
					},
				},
			},
		},
		"commentedField": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:      newIdentifier("foo"),
					Hide:    ast.ObjectFieldHidden,
					Kind:    ast.ObjectFieldID,
					Expr2:   &ast.Object{},
					Comment: &ast.Comment{Text: "a comment"},
				},
			},
		},
		"function": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Kind: ast.ObjectFieldID,
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{},
						},
					},
					Expr2: &ast.Object{},
				},
			},
		},
		"functionWithArgs": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Kind: ast.ObjectFieldID,
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{
								*newIdentifier("arg1"),
							},
						},
					},
					Expr2: &ast.Object{},
				},
			},
		},
		"binaryOp": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldHidden,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.Binary{
						Left:  &ast.Var{Id: *newIdentifier("alpha")},
						Op:    ast.BopPlus,
						Right: &ast.Var{Id: *newIdentifier("beta")},
					},
				},
			},
		},
		"conditional": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldHidden,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.Conditional{
						Cond: &ast.Binary{
							Left:  &ast.Var{Id: *newIdentifier("alpha")},
							Op:    ast.BopManifestEqual,
							Right: &ast.Var{Id: *newIdentifier("beta")},
						},
						BranchTrue: &ast.Object{
							Fields: ast.ObjectFields{
								{
									Kind: ast.ObjectFieldID,
									Hide: ast.ObjectFieldInherit,
									Id:   newIdentifier("foo"),
									Expr2: &ast.LiteralString{
										Kind:  ast.StringDouble,
										Value: "1",
									},
								},
							},
						},
						BranchFalse: &ast.Object{
							Fields: ast.ObjectFields{
								{
									Kind: ast.ObjectFieldID,
									Hide: ast.ObjectFieldInherit,
									Id:   newIdentifier("foo"),
									Expr2: &ast.LiteralString{
										Kind:  ast.StringDouble,
										Value: "2",
									},
								},
							},
						},
					},
				},
			},
		},
		"conditionalNoFalse": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Hide: ast.ObjectFieldHidden,
					Kind: ast.ObjectFieldID,
					Expr2: &ast.Conditional{
						Cond: &ast.Binary{
							Left:  &ast.Var{Id: *newIdentifier("alpha")},
							Op:    ast.BopManifestEqual,
							Right: &ast.Var{Id: *newIdentifier("beta")},
						},
						BranchTrue: &ast.Object{
							Fields: ast.ObjectFields{
								{
									Kind: ast.ObjectFieldID,
									Hide: ast.ObjectFieldInherit,
									Id:   newIdentifier("foo"),
									Expr2: &ast.LiteralString{
										Kind:  ast.StringDouble,
										Value: "1",
									},
								},
							},
						},
					},
				},
			},
		},
		"localApply": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Kind: ast.ObjectLocal,
					Expr2: &ast.Apply{
						Target: &ast.Var{
							Id: *newIdentifier("alpha"),
						},
						Arguments: ast.Arguments{
							Positional: ast.Nodes{
								&ast.LiteralString{
									Kind:  ast.StringDouble,
									Value: "arg1",
								},
							},
						},
					},
				},
			},
		},
		"localApply2": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Id:   newIdentifier("foo"),
					Kind: ast.ObjectLocal,
					Expr2: &ast.Apply{
						Target: &ast.Index{
							Id: newIdentifier("alpha"),
							Target: &ast.Index{
								Id: newIdentifier("beta"),
								Target: &ast.Var{
									Id: *newIdentifier("charlie"),
								},
							},
						},
						Arguments: ast.Arguments{
							Positional: ast.Nodes{
								&ast.LiteralString{
									Kind:  ast.StringDouble,
									Value: "arg1",
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestObject_Get(t *testing.T) {
	o := NewObject()
	v := NewObject()

	o.Set(NewKey("item"), v)

	if expected, got := v, o.Get("item"); got != expected {
		t.Fatalf("Get() got = %#v; expected = %#v", got, expected)
	}

	noder := o.Get("missing")
	if noder != nil {
		t.Fatalf("Get() nonexistant key should return nil")
	}
}

func TestBinary_UnknownOperator(t *testing.T) {
	left := NewNumber(1)
	right := NewNumber(2)

	b := NewBinary(left, right, BinaryOp("☃"))

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected unknown binary operator to panic")
		}
	}()

	b.Node()
}

func TestObject_HasUniqueKeys(t *testing.T) {
	o := NewObject()

	var err error
	err = o.Set(NewKey("foo"), NewStringDouble("text"))
	if err != nil {
		t.Errorf("Set() returned unexpected error: %v", err)
	}

	err = o.Set(NewKey("foo"), NewStringDouble("text"))
	if err == nil {
		t.Errorf("Set() expected error and there as now")
	}
}
