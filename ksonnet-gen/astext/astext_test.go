package astext

import (
	"testing"

	"github.com/google/go-jsonnet/ast"
	"github.com/stretchr/testify/require"
)

func TestCreateField(t *testing.T) {
	id := ast.Identifier("name")
	uID := ast.Identifier("underscore_name")
	leadingID := ast.Identifier("__leading")

	cases := []struct {
		name     string
		isErr    bool
		expected *ObjectField
	}{
		{
			name: "name",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldID, Id: &id}},
		},
		{
			name: "underscore_name",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldID, Id: &uID}},
		},
		{
			name: "underscore_field-",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldStr,
					Expr1: &ast.LiteralString{
						Value: "underscore_field-",
						Kind:  ast.StringDouble,
					}}},
		},
		{
			name: "dashed-name",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldStr,
					Expr1: &ast.LiteralString{
						Value: "dashed-name",
						Kind:  ast.StringDouble,
					}}},
		},
		{
			name: "__leading",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldID,
					Id:   &leadingID,
				}},
		},
		{
			name: "dot.name",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldStr,
					Expr1: &ast.LiteralString{
						Value: "dot.name",
						Kind:  ast.StringDouble,
					}}},
		},
		{
			name: ".",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldStr,
					Expr1: &ast.LiteralString{
						Value: ".",
						Kind:  ast.StringDouble,
					}}},
		},
		{
			name: "9p",
			expected: &ObjectField{
				ObjectField: ast.ObjectField{
					Kind: ast.ObjectFieldStr,
					Expr1: &ast.LiteralString{
						Value: "9p",
						Kind:  ast.StringDouble,
					}}},
		},
		{
			name:  "invalid$",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CreateField(tc.name)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, got)
			}
		})
	}
}
