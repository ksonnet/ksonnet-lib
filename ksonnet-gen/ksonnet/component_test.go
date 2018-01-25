package ksonnet_test

import (
	"path/filepath"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
	"github.com/stretchr/testify/require"
)

func testdata(name string) string {
	return filepath.Join("testdata", name)
}

func TestComponent(t *testing.T) {
	cases := []struct {
		name     string
		expected *ksonnet.Component
	}{
		{
			name: "io.k8s.api.apps.v1beta2.Deployment",
			expected: &ksonnet.Component{
				Group:   "apps",
				Version: "v1beta2",
				Kind:    "Deployment",
			},
		},
		{
			name: "io.k8s.apimachinery.pkg.apis.meta.v1.ObjectMeta",
		},
		{
			name: "missing",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			apiSpec, err := kubespec.Import(testdata("deployment.json"))
			require.NoError(t, err)

			schema := apiSpec.Definitions[tc.name]

			c := ksonnet.NewComponent(schema)

			require.Equal(t, tc.expected, c)
		})
	}

}
