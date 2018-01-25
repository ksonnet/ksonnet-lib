package kubespec_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"

	"github.com/stretchr/testify/require"
)

func testdata(name string) string {
	return filepath.Join("testdata", name)
}

func TestImporter_Import(t *testing.T) {
	cases := []struct {
		name     string
		location string
		isErr    bool
	}{
		{
			name:     "missing file",
			location: "missing.json",
			isErr:    true,
		},
		{
			name:     "invalid file",
			location: testdata("invalid.json"),
			isErr:    true,
		},
		{
			name:     "valid file",
			location: testdata("deployment.json"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Logf("path = %s", r.URL.Path)
				if r.URL.Path != "/swagger.json" {
					http.NotFound(w, r)
					return
				}

				fmt.Fprintln(w, `{"swagger": "2.0", "info": {"title": "Kubernetes"}}`)
			}))
			defer ts.Close()

			apiSpec, err := kubespec.Import(tc.location)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, apiSpec)
			}
		})
	}
}
