// Package kubeversion contains a collection of helper methods that
// help to customize the code generated for ksonnet-lib to suit
// different Kubernetes versions.
//
// For example, we may choose not to emit certain properties for some
// objects in Kubernetes v1.7.0; or, we might want to rename a
// property method. This package contains both the helper methods that
// perform such transformations, as well as the data for the
// transformations we use for each version.
package kubeversion

import (
	"log"

	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/kubespec"
)

// MapIdentifier takes a text identifier and maps it to a
// Jsonnet-appropriate identifier, for some version of Kubernetes. For
// example, in Kubernetes v1.7.0, we might map `clusterIP` ->
// `clusterIp`.
func MapIdentifier(k8sVersion, id string) string {
	verData, ok := versions[k8sVersion]
	if !ok {
		log.Fatalf("Unrecognized Kubernetes version '%s'", k8sVersion)
	}

	if alias, ok := verData.idAliases[id]; ok {
		return alias
	}
	return id
}

// IsBlacklistedProperty taks a definition name (e.g.,
// `io.k8s.kubernetes.pkg.apis.apps.v1beta1.Deployment`), a property
// name (e.g., `status`), and reports whether it is blacklisted for
// some Kubernetes version. This is particularly useful when deciding
// whether or not to generate mixins and property methods for a given
// property (as we likely wouldn't in the case of, say, `status`).
func IsBlacklistedProperty(
	k8sVersion string, path kubespec.DefinitionName,
	propertyName kubespec.PropertyName,
) bool {
	verData, ok := versions[k8sVersion]
	if !ok {
		return false
	}

	bl, ok := verData.propertyBlacklist[string(path)]
	if !ok {
		return false
	}

	_, ok = bl[string(propertyName)]
	return ok
}

//-----------------------------------------------------------------------------
// Core data structures for specifying version information.
//-----------------------------------------------------------------------------

type versionSet map[string]versionData

type versionData struct {
	idAliases         idAliasSet
	propertyBlacklist blackList
}

type idAliasSet map[string]string
type propertySet map[string]bool
type blackList map[string]propertySet

func newPropertySet(strings ...string) propertySet {
	ps := make(propertySet)
	for _, s := range strings {
		ps[s] = true
	}

	return ps
}
