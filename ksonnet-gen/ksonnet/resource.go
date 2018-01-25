package ksonnet

// Resource is a Kubernetes kind.
type Resource struct {
	description string
	properties  map[string]Field
	component   Component
	group       string
	identifier  string
}

var _ Object = (*Resource)(nil)

// NewResource creates an instance of Kind.
func NewResource(identifier, description, group string, component Component, props map[string]Field) Resource {
	return Resource{
		description: description,
		group:       group,
		component:   component,
		properties:  props,
		identifier:  identifier,
	}
}

// Kind is the kind for this resource
func (r *Resource) Kind() string {
	return r.component.Kind
}

// Version is the version for this resource
func (r *Resource) Version() string {
	return r.component.Version
}

// Group is the group for this resource
func (r *Resource) Group() string {
	if r.group == "" {
		return "core"
	}

	return r.group
}

// QualifiedGroup is the group for this resource
func (r *Resource) QualifiedGroup() string {
	return r.component.Group
}

// Description is description for this resource
func (r *Resource) Description() string {
	return r.description
}

// Identifier is identifier for this resource
func (r *Resource) Identifier() string {
	return r.identifier
}

// IsResource returns if this item is a resource. It always returns true.
func (r *Resource) IsResource() bool {
	return true
}

// Properties are the properties for this resource.
func (r *Resource) Properties() map[string]Field {
	return r.properties
}
