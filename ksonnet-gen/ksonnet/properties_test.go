package ksonnet

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_extractProperties_nil_catalog(t *testing.T) {
	_, err := extractProperties(nil, nil)
	require.Error(t, err)
}

func Test_extraProperties_nil_properties(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	m, err := extractProperties(c, nil)
	require.NoError(t, err)
	require.NotNil(t, m)
}

func Test_extraProperties_literal(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	s, ok := c.apiSpec.Definitions["io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta"]
	require.True(t, ok)

	props, err := extractProperties(c, s.Properties)
	require.NoError(t, err)

	i, ok := props["clusterName"]
	require.True(t, ok)

	prop, ok := i.(*LiteralField)
	require.True(t, ok)

	assert.Equal(t, "string", prop.FieldType())
	assert.Equal(t, "The name of the cluster which the object belongs to. This is used to distinguish resources with same name and namespace in different clusters. This field is not set anywhere right now and apiserver is going to ignore it if set in create or update request.", prop.Description())
	assert.Equal(t, "", prop.Ref())
	assert.Equal(t, "clusterName", prop.Name())
}

func Test_extraProperties_type_ref(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	s, ok := c.apiSpec.Definitions["io.k8s.api.apps.v1beta2.RollingUpdateDeployment"]
	require.True(t, ok)

	props, err := extractProperties(c, s.Properties)
	require.NoError(t, err)

	i, ok := props["maxSurge"]
	require.True(t, ok)

	prop, ok := i.(*LiteralField)
	require.True(t, ok)

	assert.Equal(t, "string", prop.FieldType())
	assert.Equal(t, "The maximum number of pods that can be scheduled above the desired number of pods. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up. Defaults to 25%. Example: when this is set to 30%, the new RC can be scaled up immediately when the rolling update starts, such that the total number of old and new pods do not exceed 130% of desired pods. Once old pods have been killed, new RC can be scaled up further, ensuring that total number of pods running at any time during the update is atmost 130% of desired pods.", prop.Description())
	assert.Equal(t, "", prop.Ref())
	assert.Equal(t, "maxSurge", prop.Name())
}

func Test_extraProperties_ref(t *testing.T) {
	c := initCatalog(t, "deployment.json")

	s, ok := c.apiSpec.Definitions["io.k8s.api.apps.v1beta2.Deployment"]
	require.True(t, ok)

	props, err := extractProperties(c, s.Properties)
	require.NoError(t, err)

	i, ok := props["metadata"]
	require.True(t, ok)

	prop, ok := i.(*ReferenceField)
	require.True(t, ok)

	assert.Equal(t, "Standard object metadata.", prop.Description())
	assert.Equal(t, "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta", prop.Ref())
	assert.Equal(t, "metadata", prop.Name())
}

func Test_extraProperties_invalid_format_ref(t *testing.T) {
	c := initCatalog(t, "invalid_ref.json")

	s, ok := c.apiSpec.Definitions["io.k8s.api.apps.v1beta2.RollingUpdateDeployment"]
	require.True(t, ok)

	_, err := extractProperties(c, s.Properties)
	require.Error(t, err)
}

func Test_fieldType(t *testing.T) {

	var (
		s1 = spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"string"},
			},
		}

		s2 = spec.Schema{}
	)

	cases := []struct {
		name     string
		schema   spec.Schema
		expected string
	}{
		{
			name:     "with an item",
			schema:   s1,
			expected: "string",
		},
		{
			name:   "with no items",
			schema: s2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fieldType(tc.schema)
			require.Equal(t, tc.expected, got)
		})
	}
}
