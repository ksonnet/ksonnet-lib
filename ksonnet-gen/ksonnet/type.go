package ksonnet

// Type is a Kubernetes type.
type Type struct {
	kind        string
	description string
	properties  map[string]Field
	version     string
	group       string
	identifier  string
}

var _ Object = (*Type)(nil)

// NewType creates an instance of Type.
func NewType(id, desc, group, ver, kind string, props map[string]Field) *Type {
	return &Type{
		identifier:  id,
		description: desc,
		group:       group,
		version:     ver,
		kind:        kind,
		properties:  props,
	}
}

// Kind is the kind for this type.
func (t *Type) Kind() string {
	return t.kind
}

// Version is the version for this type.
func (t *Type) Version() string {
	return t.version
}

// Group is the group for this type.
func (t *Type) Group() string {
	if t.group == "" {
		return "core"
	}

	return t.group
}

// QualifiedGroup is the group for this type.
func (t *Type) QualifiedGroup() string {
	return t.Group()
}

// Description is the description for this type.
func (t *Type) Description() string {
	return t.description
}

// Identifier is the identifier for this type.
func (t *Type) Identifier() string {
	return t.identifier
}

// IsResource returns if this item is a resource. It always returns false.
func (t *Type) IsResource() bool {
	return false
}

// Properties are the properties for this type.
func (t *Type) Properties() map[string]Field {
	return t.properties
}
