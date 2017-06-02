package kubespec

import (
	"fmt"
	"log"
	"strings"
)

//-----------------------------------------------------------------------------
// Utility methods for `DefinitionName`.
//-----------------------------------------------------------------------------

// Parse will parse a `DefinitionName` into a structured
// `ParsedDefinitionName`.
func (dn *DefinitionName) Parse() *ParsedDefinitionName {
	split := strings.Split(string(*dn), ".")
	if len(split) < 6 {
		log.Fatalf("Failed to parse definition name '%s'", string(*dn))
	} else if split[0] != "io" || split[1] != "k8s" || split[3] != "pkg" {
		log.Fatalf("Failed to parse definition name '%s'", string(*dn))
	}

	codebase := split[2]

	if split[4] == "api" {
		// Name is something like: `io.k8s.kubernetes.pkg.api.v1.LimitRangeSpec`.
		if len(split) < 7 {
			log.Fatalf("Failed to parse definition name '%s'", string(*dn))
		}
		versionString := VersionString(split[5])
		return &ParsedDefinitionName{
			PackageType: Core,
			Codebase:    codebase,
			Group:       nil,
			Version:     &versionString,
			Kind:        ObjectKind(split[6]),
		}
	} else if split[4] == "apis" {
		// Name is something like: `io.k8s.kubernetes.pkg.apis.batch.v1.JobList`.
		if len(split) < 8 {
			log.Fatalf("Failed to parse definition name '%s'", string(*dn))
		}
		groupName := GroupName(split[5])
		versionString := VersionString(split[6])
		return &ParsedDefinitionName{
			PackageType: APIs,
			Codebase:    codebase,
			Group:       &groupName,
			Version:     &versionString,
			Kind:        ObjectKind(split[7]),
		}
	} else if split[4] == "util" {
		if len(split) < 7 {
			log.Fatalf("Failed to parse definition name '%s'", string(*dn))
		}
		versionString := VersionString(split[5])
		return &ParsedDefinitionName{
			PackageType: Util,
			Codebase:    codebase,
			Group:       nil,
			Version:     &versionString,
			Kind:        ObjectKind(split[6]),
		}
	} else if split[4] == "runtime" {
		// Name is something like: `io.k8s.apimachinery.pkg.runtime.RawExtension`.
		return &ParsedDefinitionName{
			PackageType: Runtime,
			Codebase:    codebase,
			Group:       nil,
			Version:     nil,
			Kind:        ObjectKind(split[5]),
		}
	} else if split[4] == "version" {
		// Name is something like: `io.k8s.apimachinery.pkg.version.Info`.
		return &ParsedDefinitionName{
			PackageType: Version,
			Codebase:    codebase,
			Group:       nil,
			Version:     nil,
			Kind:        ObjectKind(split[5]),
		}
	}

	log.Fatalf("Failed to parse definition name '%s'", string(*dn))
	return nil
}

//-----------------------------------------------------------------------------
// Parsed definition name.
//-----------------------------------------------------------------------------

// Package represents the type of the definition, either `APIs`, which
// have API groups (e.g., extensions, apps, meta, and so on), or
// `Core`, which does not.
type Package int

const (
	// Core is a package that contains the Kubernetes Core objects.
	Core Package = iota

	// APIs is a set of non-core packages grouped loosely by semantic
	// functionality (e.g., apps, extensions, and so on).
	APIs

	//
	// Internal packages.
	//

	// Util is a package that contains utilities used for both testing
	// and running Kubernetes.
	Util

	// Runtime is a package that contains various utilities used in the
	// Kubernetes runtime.
	Runtime

	// Version is a package that supplies version information collected
	// at build time.
	Version
)

// ParsedDefinitionName is a parsed version of a fully-qualified
// OpenAPI spec name. For example,
// `io.k8s.kubernetes.pkg.api.v1.Container` would parse into an
// instance of the struct below.
type ParsedDefinitionName struct {
	PackageType Package
	Codebase    string
	Group       *GroupName
	Version     *VersionString
	Kind        ObjectKind
}

// GroupName represetents a Kubernetes group name (e.g., apps,
// extensions, etc.)
type GroupName string

// ObjectKind represents the `kind` of a Kubernetes API object (e.g.,
// Service, Deployment, etc.)
type ObjectKind string

// VersionString is the string representation of an API version (e.g.,
// v1, v1beta1, etc.)
type VersionString string

// Unparse transforms a `ParsedDefinitionName` back into its
// corresponding string, e.g.,
// `io.k8s.kubernetes.pkg.api.v1.Container`.
func (p *ParsedDefinitionName) Unparse() DefinitionName {
	switch p.PackageType {
	case Core:
		{
			return DefinitionName(fmt.Sprintf(
				"io.k8s.%s.pkg.api.%s.%s",
				p.Codebase,
				*p.Version,
				p.Kind))
		}
	case Util:
		{
			return DefinitionName(fmt.Sprintf(
				"io.k8s.%s.pkg.util.%s.%s",
				p.Codebase,
				*p.Version,
				p.Kind))
		}
	case APIs:
		{
			return DefinitionName(fmt.Sprintf(
				"io.k8s.%s.pkg.apis.%s.%s.%s",
				p.Codebase,
				*p.Group,
				*p.Version,
				p.Kind))
		}
	case Version:
		{
			return DefinitionName(fmt.Sprintf(
				"io.k8s.%s.pkg.version.%s",
				p.Codebase,
				p.Kind))
		}
	case Runtime:
		{
			return DefinitionName(fmt.Sprintf(
				"io.k8s.%s.pkg.runtime.%s",
				p.Codebase,
				p.Kind))
		}
	default:
		{
			log.Fatalf(
				"Failed to unparse definition name, did not recognize kind '%d'",
				p.PackageType)
			return ""
		}
	}
}
