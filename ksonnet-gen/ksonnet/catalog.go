package ksonnet

import (
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

var (
	blockedReferences = []string{
		"io.k8s.apimachinery.pkg.apis.meta.v1.ListMeta",
		"io.k8s.apiextensions-apiserver.pkg.apis.apiextensions.v1beta1.JSONSchemaProps",
		"io.k8s.apimachinery.pkg.apis.meta.v1.Status",
	}

	blockedPropertyNames = []string{
		"status",
		"$ref",
		"$schema",
		"JSONSchemas",
		"apiVersion",
		"kind",
	}
)

// ExtractFn is a function which extracts properties from a schema.
type ExtractFn func(*Catalog, map[string]spec.Schema) (map[string]Property, error)

// CatalogOpt is an option for configuring Catalog.
type CatalogOpt func(*Catalog)

// CatalogOptExtractProperties is a Catalog option for setting the property
// extractor.
func CatalogOptExtractProperties(fn ExtractFn) CatalogOpt {
	return func(c *Catalog) {
		c.extractFn = fn
	}
}

// Catalog is a catalog definitions
type Catalog struct {
	apiSpec   *spec.Swagger
	extractFn ExtractFn
}

// NewCatalog creates an instance of Catalog.
func NewCatalog(apiSpec *spec.Swagger, opts ...CatalogOpt) (*Catalog, error) {
	if apiSpec == nil {
		return nil, errors.New("apiSpec is nil")
	}

	c := &Catalog{
		apiSpec:   apiSpec,
		extractFn: extractProperties,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// Types returns a slice of all types.
func (c *Catalog) Types() ([]Type, error) {
	var resources []Type

	for name, schema := range c.definitions() {
		desc, err := ParseDescription(name)
		if err != nil {
			return nil, errors.Wrapf(err, "parse description for %s", name)
		}

		component, err := NewComponent(schema)
		if err != nil {
			continue
		}

		props, err := c.extractFn(c, schema.Properties)
		if err != nil {
			return nil, errors.Wrapf(err, "extract propererties from %s", name)
		}

		kind := NewType(name, schema.Description, desc.Group, *component, props)

		resources = append(resources, kind)
	}

	return resources, nil
}

// Fields returns a slice of all fields.
func (c *Catalog) Fields() ([]Field, error) {
	var types []Field

	for name, schema := range c.definitions() {
		desc, err := ParseDescription(name)
		if err != nil {
			return nil, errors.Wrapf(err, "parse description for %s", name)
		}

		props, err := c.extractFn(c, schema.Properties)
		if err != nil {
			return nil, errors.Wrapf(err, "extract propererties from %s", name)
		}
		t := NewField(name, schema.Description, desc.Group, desc.Version, desc.Kind, props)
		types = append(types, *t)
	}

	return types, nil
}

func (c *Catalog) isFormatRef(name string) (bool, error) {
	schema, ok := c.apiSpec.Definitions[name]
	if !ok {
		return false, errors.Errorf("%s was not found", name)
	}

	if schema.Format != "" {
		return true, nil
	}

	return false, nil
}

// Type returns a type by name. If the type cannot be found, it returns an error.
func (c *Catalog) Type(name string) (*Field, error) {
	types, err := c.Fields()
	if err != nil {
		return nil, err
	}

	for _, ty := range types {
		if ty.Identifier() == name {
			return &ty, nil
		}
	}

	return nil, errors.Errorf("%s was not found", name)
}

// Resource returns a resource by group, version, kind. If the field cannot be found,
// it returns an error
func (c *Catalog) Resource(group, version, kind string) (*Type, error) {
	resources, err := c.Types()
	if err != nil {
		return nil, err
	}

	for _, resource := range resources {
		if group == resource.Group() &&
			version == resource.Version() &&
			kind == resource.Kind() {
			return &resource, nil
		}
	}

	return nil, errors.Errorf("unable to find %s.%s.%s",
		group, version, kind)
}

func isValidDefinition(name string) bool {
	return !strings.HasPrefix(name, "io.k8s.kubernetes.pkg.api")
}

// extractRef extracts a ref from a schema.
func extractRef(schema spec.Schema) string {
	return strings.TrimPrefix(schema.Ref.String(), "#/definitions/")
}

func (c *Catalog) definitions() spec.Definitions {
	out := spec.Definitions{}

	for name, schema := range c.apiSpec.Definitions {
		if isValidDefinition(name) {
			out[name] = schema
		}
	}

	return out
}
