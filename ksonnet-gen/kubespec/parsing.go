package kubespec

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

//-----------------------------------------------------------------------------
// Utility methods for `DefinitionName` and `ObjectRef`.
//-----------------------------------------------------------------------------

// Parse will parse a `DefinitionName` into a structured
// `ParsedDefinitionName`.
func (dn *DefinitionName) Parse() *ParsedDefinitionName {
	str := string(*dn)
	packageTypeMap := map[string]Package{
		"core":    Core,
		"apis":    APIs,
		"util":    Util,
		"runtime": Runtime,
		"version": Version,
	}
	regexes := []regexp.Regexp{
		*regexp.MustCompile(`io\.k8s\.(?P<codebase>\S+)\.pkg\.api\.(?P<version>\S+)\.(?P<kind>\S+)`),                            // Core API, pre-1.8 Kubernetes OR non-Kubernetes codebase APIs
		*regexp.MustCompile(`io\.k8s\.api\.(?P<packageType>core)\.(?P<version>\S+)\.(?P<kind>\S+)`),                                    // Core API, 1.8+ Kubernetes
		*regexp.MustCompile(`io\.k8s\.(?P<codebase>\S+)\.pkg\.(?P<packageType>apis)\.(?P<group>\S+)\.(?P<version>\S+)\.(?P<kind>\S+)`), // Other APIs, pre-1.8 Kubernetes OR non-Kubernetes codebase APIs
		*regexp.MustCompile(`io\.k8s\.api\.(?P<group>\S+)\.(?P<version>\S+)\.(?P<kind>\S+)`),                                           // Other APIs, 1.8+ Kubernetes
		*regexp.MustCompile(`io\.k8s\.(?P<codebase>\S+)\.pkg\.(?P<packageType>util)\.(?P<version>\S+)\.(?P<kind>\S+)`),                 // Util packageType
		*regexp.MustCompile(`io\.k8s\.(?P<codebase>\S+)\.pkg\.(?P<packageType>runtime)\.(?P<kind>\S+)`),                                // Runtime packageType
		*regexp.MustCompile(`io\.k8s\.(?P<codebase>\S+)\.pkg\.(?P<packageType>version)\.(?P<kind>\S+)`),                                // Version packageType
	}
	for _, r := range regexes {
		if match := r.FindStringSubmatch(str); len(match) > 0 {
			result := make(map[string]string)
			for i, name := range r.SubexpNames() {
				if i != 0 {
					result[name] = match[i]
				}
			}

			// Hacky heuristics to fix missing fields
			if result["codebase"] == "" {
				result["codebase"] = "kubernetes"
			}
			if result["packageType"] == "" && result["group"] == "" {
				result["packageType"] = "core"
			}
			if result["packageType"] == "" && result["group"] != "" {
				result["packageType"] = "apis"
			}

			// Now set parsed values
			parsed := ParsedDefinitionName{}
			parsed.Codebase = result["codebase"]
			parsed.Kind = ObjectKind(result["kind"])
			parsed.PackageType = packageTypeMap[result["packageType"]]

			if result["group"] != "" {
				group := GroupName(result["group"])
				parsed.Group = &group
			}
			if result["version"] != "" {
				version := VersionString(result["version"])
				parsed.Version = &version
			}

			return &parsed
		}
	}
	log.Fatalf("Unknown definition name '%s'", str)
	return nil
}

// Name parses a `DefinitionName` from an `ObjectRef`. `ObjectRef`s
// that refer to a definition contain two parts: (1) a special prefix,
// and (2) a `DefinitionName`, so this function simply strips the
// prefix off.
func (or *ObjectRef) Name() *DefinitionName {
	defn := "#/definitions/"
	ref := string(*or)
	if !strings.HasPrefix(ref, defn) {
		log.Fatalln(ref)
	}
	name := DefinitionName(strings.TrimPrefix(ref, defn))
	return &name
}

func (dn DefinitionName) AsObjectRef() *ObjectRef {
	or := ObjectRef("#/definitions/" + dn)
	return &or
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
	Group       *GroupName     // Pointer because it's optional.
	Version     *VersionString // Pointer because it's optional.
	Kind        ObjectKind
}

// GroupName represetents a Kubernetes group name (e.g., apps,
// extensions, etc.)
type GroupName string

func (gn GroupName) String() string {
	return string(gn)
}

// ObjectKind represents the `kind` of a Kubernetes API object (e.g.,
// Service, Deployment, etc.)
type ObjectKind string

func (ok ObjectKind) String() string {
	return string(ok)
}

// VersionString is the string representation of an API version (e.g.,
// v1, v1beta1, etc.)
type VersionString string

func (vs VersionString) String() string {
	return string(vs)
}

func (p *ParsedDefinitionName) EqualStrings(p2 *ParsedDefinitionName) bool {
	return (p.Codebase == p2.Codebase) &&
	(p.PackageType == p2.PackageType) &&
	(p.Kind == p2.Kind) &&
	((p.Group == nil && p2.Group == nil) || *p.Group == *p2.Group) &&
	((p.Version == nil && p2.Version == nil) || *p.Version == *p2.Version)
}

// Unparse transforms a `ParsedDefinitionName` back into its
// corresponding string, e.g.,
// `io.k8s.kubernetes.pkg.api.v1.Container`.
func (p *ParsedDefinitionName) Unparse(withNewSchema bool) DefinitionName {
	k8s := "kubernetes"
	switch p.PackageType {
	case Core:
		{
			if withNewSchema && p.Codebase == k8s {
				return DefinitionName(fmt.Sprintf(
					"io.k8s.api.core.%s.%s",
					*p.Version,
					p.Kind))
			} else {
				return DefinitionName(fmt.Sprintf(
					"io.k8s.%s.pkg.api.%s.%s",
					p.Codebase,
					*p.Version,
					p.Kind))
			}
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
			if withNewSchema && p.Codebase == k8s {
				return DefinitionName(fmt.Sprintf(
					"io.k8s.api.%s.%s.%s",
					*p.Group,
					*p.Version,
					p.Kind))
			} else {
				return DefinitionName(fmt.Sprintf(
					"io.k8s.%s.pkg.apis.%s.%s.%s",
					p.Codebase,
					*p.Group,
					*p.Version,
					p.Kind))
			}
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
