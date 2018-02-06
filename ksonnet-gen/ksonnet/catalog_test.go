package ksonnet

import (
	"io/ioutil"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initCatalog(t *testing.T, file string, opts ...CatalogOpt) *Catalog {
	apiSpec, err := kubespec.Import(testdata(file))
	require.NoError(t, err)

	c, err := NewCatalog(apiSpec, opts...)
	require.NoError(t, err)

	return c
}

func TestCatalog_nil_apiSpec(t *testing.T) {
	_, err := NewCatalog(nil)
	require.Error(t, err)
}

func TestCatalog_Resources(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	resources, err := c.Types()
	require.NoError(t, err)

	require.Len(t, resources, 2)
}

func TestCatalog_Resources_invalid_description(t *testing.T) {
	source, err := ioutil.ReadFile("testdata/invalid_definition.json")
	require.NoError(t, err)

	apiSpec, err := kubespec.CreateAPISpec(source)
	require.NoError(t, err)

	c, err := NewCatalog(apiSpec)
	require.NoError(t, err)

	_, err = c.Types()
	assert.Error(t, err)

	_, err = c.Resource("group", "version", "kind")
	assert.Error(t, err)
}

func TestCatalog_Resources_invalid_field_properties(t *testing.T) {
	fn := func(*Catalog, map[string]spec.Schema) (map[string]Property, error) {
		return nil, errors.New("failed")
	}

	opt := CatalogOptExtractProperties(fn)

	c := initCatalog(t, "deployment.json", opt)

	_, err := c.Types()
	require.Error(t, err)
}

func TestCatalog_Resource(t *testing.T) {
	cases := []struct {
		name    string
		group   string
		version string
		kind    string
		isErr   bool
	}{
		{name: "valid id", group: "apps", version: "v1beta2", kind: "Deployment"},
		{name: "unknown kind", group: "apps", version: "v1beta2", kind: "Foo", isErr: true},
		{name: "unknown version", group: "apps", version: "Foo", kind: "Foo", isErr: true},
		{name: "unknown group", group: "Foo", version: "Foo", kind: "Foo", isErr: true},
	}

	c := initCatalog(t, "deployment.json")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r, err := c.Resource(tc.group, tc.version, tc.kind)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				t.Logf("id is %s", r.Identifier())

				require.Equal(t, tc.group, r.Group())
				require.Equal(t, tc.version, r.Version())
				require.Equal(t, tc.kind, r.Kind())
			}
		})
	}
}

func TestCatalog_Types(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	types, err := c.Fields()
	require.NoError(t, err)

	require.Len(t, types, 22)
}

func TestCatalog_Types_invalid_description(t *testing.T) {
	source, err := ioutil.ReadFile("testdata/invalid_definition.json")
	require.NoError(t, err)

	apiSpec, err := kubespec.CreateAPISpec(source)
	require.NoError(t, err)

	c, err := NewCatalog(apiSpec)
	require.NoError(t, err)

	_, err = c.Fields()
	assert.Error(t, err)

	_, err = c.Type("anything")
	assert.Error(t, err)
}

func TestCatalog_Types_invalid_field_properties(t *testing.T) {
	fn := func(*Catalog, map[string]spec.Schema) (map[string]Property, error) {
		return nil, errors.New("failed")
	}

	opt := CatalogOptExtractProperties(fn)

	c := initCatalog(t, "deployment.json", opt)

	_, err := c.Fields()
	require.Error(t, err)
}

func TestCatalog_Type(t *testing.T) {
	cases := []struct {
		name  string
		id    string
		isErr bool
	}{
		{name: "valid id", id: "io.k8s.apimachinery.pkg.apis.meta.v1.Initializers"},
		{name: "missing", id: "missing", isErr: true},
	}

	c := initCatalog(t, "deployment.json")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ty, err := c.Type(tc.id)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tc.id, ty.Identifier())
			}
		})
	}
}

func TestCatalog_isFormatRef(t *testing.T) {
	cases := []struct {
		name        string
		isFormatRef bool
		isErr       bool
	}{
		{
			name: "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		},
		{
			name:  "missing",
			isErr: true,
		},
		{
			name:        "io.k8s.apimachinery.pkg.util.intstr.IntOrString",
			isFormatRef: true,
		},
	}

	c := initCatalog(t, "deployment.json")

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tf, err := c.isFormatRef(tc.name)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				require.Equal(t, tc.isFormatRef, tf)
			}
		})
	}
}
