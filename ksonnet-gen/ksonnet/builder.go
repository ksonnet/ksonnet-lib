package ksonnet

import (
	"log"

	k8s "github.com/ksonnet/ksonnet-lib/ksonnet-gen/k8sSwagger"
)

type builder interface {
	Root() *namespace
	AddAppSpec(model *k8s.AppSpec)
}

func newBuilder(apiVersion string) builder {
	return &builderCore{
		rootNamespace: makeRootNamespace(apiVersion),
	}
}

type builderCore struct {
	rootNamespace *namespace
}

func (b *builderCore) Root() *namespace {
	return b.rootNamespace
}

func (b *builderCore) AddAppSpec(appSpec *k8s.AppSpec) {
	for _, model := range appSpec.Models {
		b.addModel(model)
	}

	b.createMixins(b.Root())
}

func (b *builderCore) addModel(model k8s.Model) {
	// Recursively create a namespace for every component of the path
	// (if they don't exist already). For example, `v1.foo.bar` creates
	// a namespace `v1`, which contains a namespace `foo`, and so on.
	ns, _ := b.getOrCreateNamespace(model.ID)

	ns.constructor = makeConstructor(ns, model)

	// Create a function for each property that is optional.
	for propertyName, property := range model.Properties {
		if propertyName == constructorName {
			// log.Fatalf(
			// 	"Properties called `default` are disallowed: in namespace '%s'",
			// 	ns.ID)
			id := propertyName + "Value"
			ns.addPropertyMethod(id, id, propertyName, property)
		} else if propertyName == errorKeyword {
			id := propertyName + "Condition"
			ns.addPropertyMethod(id, id, propertyName, property)
		} else {
			ns.addPropertyMethod(propertyName, propertyName, propertyName, property)
		}
	}
}

func (b *builderCore) createMixins(currNS *namespace) {
	metadataMethod := currNS.specialMethods[objectMetadata]
	b.createMixinSpecFromMethod(currNS, metadataMethod)
	for _, method := range currNS.methods {
		b.createMixinSpecFromMethod(currNS, method)
	}

	for _, childNS := range currNS.namespaceSet {
		b.createMixins(childNS)
	}
}

func (b *builderCore) createMixinSpecFromMethod(
	currNS *namespace, method *propertyMethod,
) {
	//
	// Get namespace for the method ref, if it exists. This should be
	// something like `v1.ObjectMeta`. So we need to go get all the
	// methods in the object that ref points at, and create a mixin for
	// each.
	//

	if method == nil || method.Ref == nil {
		return
	}
	refNS, created := b.getOrCreateNamespace(*method.Ref)
	if created {
		log.Fatalf(
			"Method '%s.%s' references namespace '%s', but that path does not exist",
			currNS.ID,
			method.name,
			*method.Ref)
	}

	currNS.mixins[method.fieldName] = newMixinSpec(
		method.fieldName, method.fieldName, refNS, currNS)
}

func (b *builderCore) getOrCreateNamespace(path string) (*namespace, bool) {
	created := false
	nsPath := parsePath(path)
	// Create or add to the definition of each namespace in the path,
	// where the path is something like `v1.foo.bar`.
	currNS := b.Root()
	for i, nsID := range nsPath {
		if nsObject, nsExists := currNS.namespaceSet[nsID]; nsExists {
			currNS.namespaceSet[nsID] = nsObject
		} else {
			created = true
			currPath := unparsePath(nsPath[:i+1])
			currNS.namespaceSet[nsID] = makeNamespace(currPath, nsID)
		}
		currNS = currNS.namespaceSet[nsID]
	}

	return currNS, created
}
