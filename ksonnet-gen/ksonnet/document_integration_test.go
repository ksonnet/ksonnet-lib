package ksonnet_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/ksonnet"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/stretchr/testify/require"
)

func testdata(name string) string {
	return filepath.Join("testdata", name)
}

func TestDocument_Integration(t *testing.T) {
	dir, err := ioutil.TempDir("", "document")
	require.NoError(t, err)

	defer os.RemoveAll(dir)

	b := genDoc(t, "swagger-1.8.json")

	k8sPath := filepath.Join(dir, "k8s.libsonnet")
	writeFile(t, k8sPath, b)

	verifyK8s(t, dir)

	ksPath := filepath.Join(dir, "k.libsonnet")
	copyFile(t, testdata("k.libsonnet"), ksPath)

	compPath := filepath.Join(dir, "component.libsonnet")
	copyFile(t, testdata("component.libsonnet"), compPath)

	cmd := exec.Command(jsonnetCmd(), "component.libsonnet")
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out))
	}

	expected, err := ioutil.ReadFile(testdata("component.json"))
	require.NoError(t, err)

	require.Equal(t, string(expected), string(out))
}

func jsonnetCmd() string {
	bin := os.Getenv("JSONNET_BIN")
	if bin == "" {
		bin = "jsonnet"
	}

	return bin
}

func verifyK8s(t *testing.T, dir string) {
	cmd := exec.Command(jsonnetCmd(), "fmt", "k8s.libsonnet")
	cmd.Dir = dir

	var b bytes.Buffer
	cmd.Stderr = &b

	err := cmd.Run()
	if err != nil {
		t.Fatalf("k8s.libsonnet verification failed: %v", b.String())
	}
}

func genDoc(t *testing.T, input string) []byte {
	apiSpec, checksum, err := kubespec.Import(testdata(input))
	require.NoError(t, err)

	c, err := ksonnet.NewCatalog(apiSpec, ksonnet.CatalogOptChecksum(checksum))
	require.NoError(t, err)

	doc, err := ksonnet.NewDocument(c)
	require.NoError(t, err)

	node, err := doc.Node()
	require.NoError(t, err)

	var buf bytes.Buffer
	require.NoError(t, printer.Fprint(&buf, node.Node()))

	return buf.Bytes()
}

func writeFile(t *testing.T, name string, content []byte) {
	err := ioutil.WriteFile(name, content, 0600)
	require.NoError(t, err)
}

func copyFile(t *testing.T, src, dest string) {
	from, err := os.Open(src)
	require.NoError(t, err)
	defer from.Close()

	to, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
	require.NoError(t, err)
	defer to.Close()

	_, err = io.Copy(to, from)
	require.NoError(t, err)
}
