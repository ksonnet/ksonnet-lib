package ksonnet

import (
	nm "github.com/ksonnet/ksonnet-lib/ksonnet-gen/nodemaker"
	"github.com/pkg/errors"
)

// APIObject is an API object.
type APIObject struct {
	resource       Object
	renderFieldsFn renderFieldsFn
}

// NewAPIObject creates an instance of APIObject.
func NewAPIObject(resource Object) *APIObject {
	ao := &APIObject{
		resource:       resource,
		renderFieldsFn: renderFields,
	}

	return ao
}

// Kind is the kind of api object this is.
func (a *APIObject) Kind() string {
	return FormatKind(a.resource.Kind())
}

// Description is the description of this API object.
func (a *APIObject) Description() string {
	return a.resource.Description()
}

// Node returns an AST node for this api object.
func (a *APIObject) Node(catalog *Catalog) (*nm.Object, error) {
	return apiObjectNode(catalog, a)
}

func (a *APIObject) initNode() *nm.Object {
	o := nm.NewObject()

	if a.resource.IsType() {
		kindObject := nm.OnelineObject()
		kind := a.resource.Kind()
		kindObject.Set(nm.InheritedKey("kind"), nm.NewStringDouble(kind))
		o.Set(nm.LocalKey("kind"), kindObject)

		o.Set(nm.FunctionKey("new", []string{}), objectConstructor())
	} else {
		o.Set(nm.FunctionKey("new", []string{}), nm.OnelineObject())
	}

	return o
}

func objectConstructor() *nm.Binary {
	return nm.NewBinary(nm.NewVar("apiVersion"), nm.NewVar("kind"), nm.BopPlus)
}

func apiObjectNode(catalog *Catalog, a *APIObject) (*nm.Object, error) {
	if catalog == nil {
		return nil, errors.New("catalog is nil")
	}

	o := a.initNode()
	if err := a.renderFieldsFn(catalog, o, "", a.resource.Properties()); err != nil {
		return nil, err
	}
	return o, nil
}
