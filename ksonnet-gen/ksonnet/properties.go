package ksonnet

import (
	"strings"

	"github.com/go-openapi/spec"
	"github.com/pkg/errors"
)

func extractProperties(c *Catalog, properties map[string]spec.Schema) (map[string]Property, error) {
	if c == nil {
		return nil, errors.New("catalog is nil")
	}

	out := make(map[string]Property)

	for name, schema := range properties {
		if isSkippedProperty(name, schema) {
			continue
		}

		ref := extractRef(schema)

		// literal
		if t := schema.Type; len(t) == 1 {
			out[name] = buildLiteralField(t[0], name, schema)
			continue
		}

		ifr, err := c.isFormatRef(ref)
		if err != nil {
			return nil, errors.Wrap(err, "check for format ref")
		}

		if ifr {
			// don't have to check for existence here because isFormatRef does the same thing
			formatSchema := c.apiSpec.Definitions[ref]
			out[name] = buildLiteralField(fieldType(formatSchema), name, schema)
			continue
		}

		// must be a mixin
		f := NewReferenceField(name, schema.Description, ref)
		out[name] = f
	}

	return out, nil
}

func buildLiteralField(fieldType, name string, schema spec.Schema) *LiteralField {
	var itemRef string
	if schema.Items != nil && schema.Items.Schema != nil {
		itemRef = extractRef(*schema.Items.Schema)
	}

	return NewLiteralField(name, fieldType, schema.Description, itemRef)
}

func isSkippedProperty(name string, schema spec.Schema) bool {
	if stringInSlice(name, blockedPropertyNames) {
		return true
	}

	if strings.Contains(strings.ToLower(schema.Description), "read-only") {
		return true
	}

	ref := extractRef(schema)
	if stringInSlice(ref, blockedReferences) {
		return true
	}

	return false
}

func fieldType(schema spec.Schema) string {
	if t := schema.Type; len(t) == 1 {
		return t[0]
	}

	return ""
}
