package ksonnet

import (
	"io/ioutil"
	"testing"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/stretchr/testify/require"
)

func TestExtension_Output(t *testing.T) {
	c := initCatalog(t, "swagger-1.8.json")

	e := NewExtension(c)

	node, err := e.Node()
	require.NoError(t, err)

	require.NoError(t, printer.Fprint(ioutil.Discard, node.Node()))
}
