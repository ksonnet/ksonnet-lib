package printer

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ast"
)

func TestFprintf(t *testing.T) {
	cases := []struct {
		name  string
		isErr bool
	}{
		{name: "object"},
		{name: "object_with_hidden_field"},
		{name: "inline_object"},
		{name: "object_mixin"},
		{name: "object_with_nested_object"},
		{name: "local"},
		{name: "multi_line_comments"},
		{name: "literal"},
		{name: "literal_with_newline"},
		{name: "literal_with_single_quote"},
		{name: "object_field_with_comment"},
		{name: "function_with_no_args"},
		{name: "function_with_args"},
		{name: "function_with_optional_args"},
		{name: "local_function_with_args"},
		{name: "conditional"},
		{name: "conditional_no_false"},
		{name: "index"},
		{name: "index_with_index"},
		{name: "array"},

		// errors
		{name: "unknown_node", isErr: true},
		{name: "nil_node", isErr: true},
		{name: "invalid_apply", isErr: true},
		{name: "invalid_literal_string", isErr: true},
		{name: "invalid_of_kind", isErr: true},
		{name: "invalid_of_hide", isErr: true},
		{name: "invalid_of_method", isErr: true},
		{name: "index_no_index_or_id", isErr: true},
		{name: "index_invalid_index", isErr: true},
		{name: "index_invalid_literal_string", isErr: true},
		{name: "null_index", isErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			node, ok := fprintfCases[tc.name]
			if !ok {
				t.Fatalf("test case %q does not exist", tc.name)
			}

			err := Fprint(&buf, node)
			if tc.isErr {
				if err == nil {
					t.Fatalf("test case %q expected error and it was not", tc.name)
				}
			} else {
				if err != nil {
					t.Fatalf("test case %q returned an unexpected error: %v", tc.name, err)
				}

				testDataFile := filepath.Join("testdata", tc.name)
				testData, err := ioutil.ReadFile(testDataFile)

				if err != nil {
					t.Fatalf("unable to read test data: %v", err)
				}

				if got, expected := buf.String(), string(testData); got != expected {
					t.Fatalf("Fprint\ngot      = %s\nexpected = %s",
						strconv.Quote(got), strconv.Quote(expected))
				}
			}

		})
	}
}

var (
	id1 = ast.Identifier("foo")
	id2 = ast.Identifier("bar")

	fprintfCases = map[string]ast.Node{
		"object": &ast.Object{
			Fields: ast.ObjectFields{},
		},
		"object_with_hidden_field": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind:  ast.ObjectFieldID,
					Id:    newIdentifier("foo"),
					Expr2: &ast.Object{},
				},
				{
					Kind:  ast.ObjectFieldID,
					Id:    newIdentifier("bar"),
					Expr2: &ast.Object{},
				},
			},
		},
		"inline_object": &ast.Object{
			Oneline: true,
			Fields: ast.ObjectFields{
				{
					Kind:  ast.ObjectFieldID,
					Id:    &id1,
					Expr2: &ast.Var{Id: *newIdentifier("foo")},
				},
				{
					Kind:  ast.ObjectFieldID,
					Id:    &id1,
					Expr2: &ast.Var{Id: *newIdentifier("bar")},
				},
			},
		},
		"index": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   &id1,
					Expr2: &ast.Index{
						Id: newIdentifier("foo"),
						Target: &ast.Index{
							Id: newIdentifier("bar"),
							Target: &ast.Var{
								Id: *newIdentifier("baz"),
							},
						},
					},
				},
			},
		},
		"index_with_index": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   &id1,
					Expr2: &ast.Index{
						Id: newIdentifier("foo"),
						Target: &ast.Index{
							Index: &ast.LiteralString{
								Value: "bar",
								Kind:  ast.StringDouble,
							},
							Target: &ast.Var{
								Id: *newIdentifier("baz"),
							},
						},
					},
				},
			},
		},
		"object_mixin": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind:       ast.ObjectFieldID,
					Id:         &id1,
					Expr2:      &ast.Object{},
					SuperSugar: true,
				},
			},
		},
		"object_with_nested_object": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   &id1,
					Expr2: &ast.Object{
						Fields: ast.ObjectFields{
							{
								Kind:  ast.ObjectFieldID,
								Id:    &id2,
								Expr2: &ast.Object{},
							},
						},
					},
				},
			},
		},
		"local": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind:  ast.ObjectLocal,
					Hide:  ast.ObjectFieldVisible,
					Id:    &id2,
					Expr2: &ast.Object{},
				},
			},
		},
		"multi_line_comments": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind:  ast.ObjectLocal,
					Hide:  ast.ObjectFieldVisible,
					Id:    &id2,
					Expr2: &ast.Object{},
					Comment: &ast.Comment{
						Text: "line 1\n\nline 3\nline 4",
					},
				},
			},
		},
		"literal": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Hide: ast.ObjectFieldInherit,
					Id:   &id1,
					Expr2: &ast.LiteralString{
						Value: "value",
						Kind:  ast.StringDouble,
					},
				},
			},
		},
		"literal_with_newline": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Hide: ast.ObjectFieldInherit,
					Id:   &id1,
					Expr2: &ast.LiteralString{
						Value: "value1\nvalue2",
						Kind:  ast.StringDouble,
					},
				},
			},
		},
		"literal_with_single_quote": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Hide: ast.ObjectFieldInherit,
					Id:   &id1,
					Expr2: &ast.LiteralString{
						Value: "value1",
						Kind:  ast.StringSingle,
					},
				},
			},
		},
		"object_field_with_comment": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Hide: ast.ObjectFieldInherit,
					Id:   &id1,
					Expr2: &ast.LiteralString{
						Value: "value",
						Kind:  ast.StringDouble,
					},
					Comment: &ast.Comment{
						Text: "a comment",
					},
				},
			},
		},
		"function_with_no_args": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   &id1,
					Expr2: &ast.Binary{
						Left:  newLiteralNumber("1"),
						Right: newLiteralNumber("2"),
						Op:    ast.BopPlus,
					},
					Method: &ast.Function{
						Parameters: ast.Parameters{},
					},
				},
			},
		},
		"function_with_args": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   &id1,
					Expr2: &ast.Binary{
						Left:  &ast.Var{Id: *newIdentifier("myVar")},
						Right: newLiteralNumber("2"),
						Op:    ast.BopPlus,
					},
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{
								*newIdentifier("one"),
								*newIdentifier("two"),
							},
						},
					},
				},
			},
		},
		"function_with_optional_args": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectFieldID,
					Id:   newIdentifier("alpha"),
					Expr2: &ast.Binary{
						Left:  &ast.Var{Id: *newIdentifier("myVar")},
						Right: newLiteralNumber("2"),
						Op:    ast.BopPlus,
					},
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{
								*newIdentifier("one"),
								*newIdentifier("two"),
							},
							Optional: []ast.NamedParameter{
								{
									Name:       *newIdentifier("opt1"),
									DefaultArg: newLiteralNumber("1"),
								},
							},
						},
					},
				},
				{
					Kind: ast.ObjectFieldID,
					Id:   newIdentifier("beta"),
					Expr2: &ast.Binary{
						Left:  &ast.Var{Id: *newIdentifier("myVar")},
						Right: newLiteralNumber("2"),
						Op:    ast.BopPlus,
					},
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{
								*newIdentifier("one"),
								*newIdentifier("two"),
							},
							Optional: []ast.NamedParameter{
								{
									Name: *newIdentifier("opt1"),
									DefaultArg: &ast.Object{
										Oneline: true,
										Fields: ast.ObjectFields{
											{
												Kind:  ast.ObjectFieldID,
												Hide:  ast.ObjectFieldInherit,
												Id:    newIdentifier("foo"),
												Expr2: newLiteralNumber("1"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		"local_function_with_args": &ast.Object{
			Fields: ast.ObjectFields{
				{
					Kind: ast.ObjectLocal,
					Id:   newIdentifier("foo"),
					Expr2: &ast.Binary{
						Left:  &ast.Var{Id: *newIdentifier("myVar")},
						Right: newLiteralNumber("2"),
						Op:    ast.BopPlus,
					},
					Method: &ast.Function{
						Parameters: ast.Parameters{
							Required: ast.Identifiers{
								*newIdentifier("one"),
								*newIdentifier("two"),
							},
						},
					},
				},
			},
		},
		"conditional": &ast.Conditional{
			Cond: &ast.Binary{
				Left: &ast.Apply{
					Target: &ast.Index{
						Id: newIdentifier("std"),
						Target: &ast.Var{
							Id: *newIdentifier("type"),
						},
					},
					Arguments: ast.Arguments{
						Positional: ast.Nodes{
							&ast.Var{Id: *newIdentifier("foo")},
						},
					},
				},
				Right: &ast.LiteralString{
					Value: "array",
					Kind:  ast.StringDouble,
				},
				Op: ast.BopManifestEqual,
			},
			BranchTrue: &ast.Object{
				Oneline: true,
				Fields: ast.ObjectFields{
					{
						Id:    newIdentifier("foo"),
						Kind:  ast.ObjectFieldID,
						Hide:  ast.ObjectFieldInherit,
						Expr2: &ast.Var{Id: *newIdentifier("foo")},
					},
				},
			},
			BranchFalse: &ast.Object{
				Oneline: true,
				Fields: ast.ObjectFields{
					{
						Id:   newIdentifier("foo"),
						Kind: ast.ObjectFieldID,
						Hide: ast.ObjectFieldInherit,
						Expr2: &ast.Array{
							Elements: ast.Nodes{
								&ast.Var{Id: *newIdentifier("foo")},
							},
						},
					},
				},
			},
		},
		"conditional_no_false": &ast.Conditional{
			Cond: &ast.Binary{
				Left: &ast.Apply{
					Target: &ast.Index{
						Id: newIdentifier("std"),
						Target: &ast.Var{
							Id: *newIdentifier("type"),
						},
					},
					Arguments: ast.Arguments{
						Positional: ast.Nodes{
							&ast.Var{Id: *newIdentifier("foo")},
						},
					},
				},
				Right: &ast.LiteralString{
					Value: "array",
					Kind:  ast.StringDouble,
				},
				Op: ast.BopManifestEqual,
			},
			BranchTrue: &ast.Object{
				Oneline: true,
				Fields: ast.ObjectFields{
					{
						Id:    newIdentifier("foo"),
						Kind:  ast.ObjectFieldID,
						Hide:  ast.ObjectFieldInherit,
						Expr2: &ast.Var{Id: *newIdentifier("foo")},
					},
				},
			},
		},
		"array": &ast.Array{
			Elements: ast.Nodes{
				&ast.Var{Id: *newIdentifier("foo")},
				&ast.Self{},
				&ast.LiteralString{
					Value: "string",
				},
			},
		},

		// errors
		"unknown_node":           &noopNode{},
		"nil_node":               nil,
		"invalid_apply":          &ast.Apply{Target: newLiteralNumber("1")},
		"invalid_literal_string": &ast.LiteralString{Kind: 99},
		"invalid_of_kind": &ast.Object{
			Fields: ast.ObjectFields{{Kind: 99}},
		},
		"invalid_of_hide": &ast.Object{
			Fields: ast.ObjectFields{{Hide: 99}},
		},
		"invalid_of_method": &ast.Object{
			Fields: ast.ObjectFields{{
				Method: &ast.Function{
					Parameters: ast.Parameters{
						Optional: []ast.NamedParameter{
							{Name: *newIdentifier("opt1"), DefaultArg: &noopNode{}},
						},
					},
				},
			}},
		},
		"index_no_index_or_id":         &ast.Index{},
		"index_invalid_index":          &ast.Index{Index: &noopNode{}},
		"index_invalid_literal_string": &ast.Index{Index: (*ast.LiteralString)(nil)},
		"null_index":                   (*ast.Index)(nil),
	}
)

func Test_printer_indent(t *testing.T) {
	cases := []struct {
		name     string
		level    int
		expected string
		output   []byte
		mode     IndentMode
	}{
		{
			name:  "space: empty",
			level: 0, expected: "", output: make([]byte, 0)},
		{
			name:  "space: not at eol",
			level: 0, expected: "word", output: []byte("word")},
		{
			name:  "space: at eol",
			level: 0, expected: "word\n", output: []byte("word\n")},
		{
			name:  "space: indent level 1: empty",
			level: 1, expected: "", output: make([]byte, 0)},
		{
			name:  "space: indent level 1: not at eol",
			level: 1, expected: "word", output: []byte("word")},
		{
			name:  "space: indent level 1: at eol",
			level: 1, expected: "word\n  ", output: []byte("word\n")},
		{
			name:  "tab: empty",
			level: 0, expected: "", output: make([]byte, 0), mode: IndentModeTab},
		{
			name:  "tab: not at eol",
			level: 0, expected: "word", output: []byte("word"), mode: IndentModeTab},
		{
			name:  "tab: at eol",
			level: 0, expected: "word\n", output: []byte("word\n"), mode: IndentModeTab},
		{
			name:  "tab: indent level 1: empty",
			level: 1, expected: "", output: make([]byte, 0), mode: IndentModeTab},
		{
			name:  "tab: indent level 1: not at eol",
			level: 1, expected: "word", output: []byte("word"), mode: IndentModeTab},
		{
			name:  "tab: indent level 1: at eol",
			level: 1, expected: "word\n\t", output: []byte("word\n"), mode: IndentModeTab},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := printer{cfg: Config{IndentMode: tc.mode, IndentSize: 2}}
			p.indentLevel = tc.level

			for _, b := range tc.output {
				p.writeByte(b, 1)
			}

			expected := tc.expected
			if got := string(p.output); got != expected {
				t.Fatalf("Fprint\ngot      = %s\nexpected = %s",
					strconv.Quote(got), strconv.Quote(expected))
			}
		})
	}
}

func Test_printer_indent_empty(t *testing.T) {
	p := printer{cfg: DefaultConfig}
	p.indentLevel = 1
	p.indent()
	if len(p.output) != 0 {
		t.Errorf("indent() with empty output should not change output")
	}
}

func Test_printer_err(t *testing.T) {
	p := printer{cfg: DefaultConfig}
	p.err = errors.New("error")

	n := &ast.Object{}
	p.print(n)

	if len(p.output) != 0 {
		t.Errorf("print() in error state should not add any output")
	}
}

func newLiteralNumber(in string) *ast.LiteralNumber {
	f, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return &ast.LiteralNumber{OriginalString: in, Value: 0}
	}
	return &ast.LiteralNumber{OriginalString: in, Value: f}
}

// newIdentifier creates an identifier.
func newIdentifier(value string) *ast.Identifier {
	id := ast.Identifier(value)
	return &id
}

type noopNode struct{}

func (n *noopNode) Context() ast.Context             { return nil }
func (n *noopNode) Loc() *ast.LocationRange          { return nil }
func (n *noopNode) FreeVariables() ast.Identifiers   { return nil }
func (n *noopNode) SetFreeVariables(ast.Identifiers) {}
func (n *noopNode) SetContext(ast.Context)           {}
