package ksonnet

import (
	"log"
	"strings"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
)

var specialProperties = map[kubespec.PropertyName]kubespec.PropertyName{
	"apiVersion": "apiVersion",
	"metadata":   "metadata",
	"kind":       "kind",
}

func isSpecialProperty(pn kubespec.PropertyName) bool {
	_, ok := specialProperties[pn]
	return ok
}

func toJsonnetName(ok kubespec.ObjectKind) kubespec.ObjectKind {
	if len(ok) == 0 {
		log.Fatalf("Can't lowercase first letter of 0-rune string")
	}
	kindString := string(ok)

	upper := strings.ToLower(kindString[:1])
	return kubespec.ObjectKind(upper + kindString[1:])
}
