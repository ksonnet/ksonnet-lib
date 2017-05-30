package ksonnet

import (
	"fmt"
	"sort"
	"strings"

	k8s "github.com/ksonnet/ksonnet-lib/ksonnet-gen/k8sSwagger"
)

//-----------------------------------------------------------------------------
// Namespace.
//-----------------------------------------------------------------------------

// namespace is a specially-structured Jsonnet object that contains
// everything we need to manage a Kubernetes API object. For example,
// for `v1.Container`, this object would contain property methods,
// like:
//
//   container.image(image)
//
// mixins, like:
//
//   container.mixin.livenessProbe.initialDelaySeconds(seconds)
//
// as well as "special" methods that set things like `"apiVersion"`,
// where applicable.
type namespace struct {
	ID             string
	path           string
	APIVersion     *string
	namespaceSet   namespaceSet
	constructor    *constructor
	methods        propertyMethodSet
	specialMethods propertyMethodSet
	mixins         mixinSpecSet
}

func makeRootNamespace(apiVersion string) *namespace {
	rootNamespaceID := "$"
	ns := makeNamespace(rootNamespaceID, rootNamespaceID)
	ns.APIVersion = &apiVersion
	return ns
}

func makeNamespace(path, id string) *namespace {
	return &namespace{
		ID:             id,
		path:           path,
		constructor:    nil,
		namespaceSet:   make(map[string]*namespace),
		methods:        makePropertyMethodSet(),
		specialMethods: makePropertyMethodSet(),
		mixins:         makeMixinSpecSet(),
	}
}

func (ns *namespace) addPropertyMethod(
	name, param, fieldName string, property k8s.Property,
) {
	//
	// Property methods, somewhat redundantly, look something like:
	//
	//   name(name):: {name: name}
	//
	// The caller doesn't see this redundancy. They typically use the
	// the mixin facilities, something like:
	// `metadata.name("foo") + metadata.namespace("bar")`.
	//
	method := newPropertyMethod(name, param, fieldName, property)

	if _, ok := specialResourcePropertiesSet[name]; ok {
		// Don't expose functions for stuff like `kind` and
		// `apiVersion`. That should be populated in the constructor,
		// for free, without the user needing to do it manually.
		ns.specialMethods[name] = &method
	} else {
		ns.methods[name] = &method
	}
}

//
// The namespaceSet type.
//

type namespaceSet map[string]*namespace

func (nsSet namespaceSet) toSlice() namespaces {
	nss := namespaces{}
	for _, of := range nsSet {
		nss = append(nss, of)
	}

	sort.Sort(nss)
	return nss
}

func (nsSet namespaceSet) addNamespace(ns *namespace) namespaceSet {
	if ns != nil {
		nsSet[ns.ID] = ns
	}
	return nsSet
}

//
// The namespaces type.
//

type namespaces []*namespace

func (nss namespaces) Len() int {
	return len(nss)
}

func (nss namespaces) Less(i, j int) bool {
	return nss[i].ID < nss[j].ID
}

func (nss namespaces) Swap(i, j int) {
	nss[i], nss[j] = nss[j], nss[i]
}

//-----------------------------------------------------------------------------
// Constructor.
//-----------------------------------------------------------------------------

// constructor represents a Kubernetes API object's constructor. For
// example, `v1.Container` might have a constructor like:
//
//   container.default(imageName)
//
// This struct represents the definition of this method. Note that
// this can be fairly complex, as we may have to add initial values to
// several fields (such as `"kind"` and `"apiVersion"`).
type constructor struct {
	parent      *namespace
	params      []constructorParam
	assignments assignmentSet
	name        string
}

func makeConstructor(parent *namespace, model k8s.Model) *constructor {
	// Initialize complex types as their default values.
	// NOTE: It's necessary for this loop to run first, as we want the
	// second to override these default values with the references to
	// the parameters.
	assignments := map[string]string{}
	for propertyName, property := range model.Properties {
		if propertyName == objectMetadata {
			// The metadata field will be set using a call to
			// `$.v1.metadata.default`.
			continue
		}

		if property.Type != nil {
			switch *property.Type {
			case "object":
				assignments[propertyName] = "{}"
			case "array":
				assignments[propertyName] = "[]"
			}
		} else if property.Ref != nil {
			assignments[propertyName] = "{}"
		}
	}

	// Create a `default` function with all required properties. (That
	// number may be 0).
	requiredParams := []constructorParam{}
	for _, propertyName := range model.Required {
		// Skip special methods like `apiVersion` and `kind`.
		if _, ok := specialResourcePropertiesSet[propertyName]; ok {
			continue
		}

		requiredParam :=
			constructorParam{propertyName, model.Properties[propertyName]}
		requiredParams = append(requiredParams, requiredParam)

		// Autobox constructor parameters where type == "array". For
		// example, something like `deployment.default` might take a
		// container or array of containers as argument. In the case of
		// the former, we just turn it into an array.
		property := model.Properties[requiredParam.id]
		var rhs string
		if property.Type != nil && (*property.Type == "array") {
			rhs = fmt.Sprintf(
				"if std.type(%s) == \"array\" then %s else [%s]",
				requiredParam.id,
				requiredParam.id,
				requiredParam.id)
		} else {
			rhs = requiredParam.id
		}

		assignments[requiredParam.id] = rhs
	}

	return &constructor{
		parent:      parent,
		params:      requiredParams,
		assignments: assignments,
		name:        constructorName,
	}
}

func (method *constructor) emit(m marshaller) error {
	//
	// Aggregate information about a method, but defer writing to the
	// buffer until we're sure it's not a no-op; if it is, don't write
	// it at all.
	//
	// First, assemble the parameter list, which goes in the function
	// signature (e.g., the `name, namespace="default"` in
	// `default(name, namespace="default")`), and the assignment list,
	// which goes in the function body (e.g., the `name: name` that
	// assigns the `name` parameter to the `name` field in the body of
	// the function).
	//

	paramList := []string{}
	for _, param := range method.params {
		paramList = append(paramList, param.id)
	}

	if _, ok := method.parent.specialMethods[objectMetadata]; method.name == constructorName && ok {
		paramList = append([]string{metadataName}, paramList...)
		paramList = append(paramList, metadataNamespace+"=\"default\"")
	}

	//
	// Second, assemble and write out the signature.
	//

	signature := fmt.Sprintf(
		"%s(%s)::", method.name, strings.Join(paramList, ", "))

	//
	// Third, write out the hard-coded values for the "special methods".
	// These are things like "kind", "apiVersion", or "metadata.name",
	// which are required for several objects, but which we can infer,
	// and therefore don't need the user to manually write them in. We
	// use Jsonnet's mixin facilities to handle this, so that it
	// generates the `+`'d values here (and not the signature or body):
	//
	//   fooProperty()::
	//     apiVersion +
	//     kind +
	//     {
	//        [...]
	//     }
	//

	var specialCalls = []string{}
	if method.name == constructorName {
		for _, specialMethod := range method.parent.specialMethods.toSlice() {
			if specialMethod.name == objectMetadata {
				// All objects that have an `apiVersion` are required to have
				// a `metadata` field with `name` and `namespace`. Here we use
				// the metadata functions to build this up.
				specialCalls = append(
					specialCalls,
					fmt.Sprintf(
						"%s(%s, %s) +",
						defaultMetadataFunctionName,
						metadataName,
						metadataNamespace))

			} else {
				specialCalls = append(
					specialCalls,
					fmt.Sprintf("%s +", specialMethod.name))
			}
		}
	}

	//
	// Write the body and the assignments. This should look something like:
	//
	//   {
	//     name: name,
	//   }
	//

	bodyOpenBrace := "{"

	assignmentStrings := []string{}
	for _, assignment := range method.assignments.toSlice() {
		assignmentStrings = append(
			assignmentStrings,
			fmt.Sprintf("%s: %s,", assignment.fieldName, assignment.fieldVal))
	}

	bodyClosingBrace := "},"

	//
	// Write out if the method isn't a no-op.
	//

	if len(paramList) == 0 && len(specialCalls) == 0 && len(method.assignments) == 0 {
		// Don't bother writing a completely empty method.
		return nil
	}

	m.bufferLine(signature)
	m.indent()

	for _, specialCall := range specialCalls {
		m.bufferLine(specialCall)
	}
	m.dedent()

	m.bufferLine(bodyOpenBrace)
	m.indent()

	for _, assignment := range assignmentStrings {
		m.bufferLine(assignment)
	}
	m.dedent()

	m.bufferLine(bodyClosingBrace)

	return nil
}

//
// The constructorParam type.
//

type constructorParam struct {
	id       string
	property k8s.Property
}

//-----------------------------------------------------------------------------
// Assignments.
//-----------------------------------------------------------------------------

// assignment represents the assignment (strictly, a definition, but
// that terminology is more confusing) of a field to some value in an
// object. For example, in the following object, `assignment` would
// represent the assignment of value 99 to property `foo`:
//
//   {"foo": 99}
//
type assignment struct {
	fieldName string
	fieldVal  string
}

//
// The assignmentSet type.
//

type assignmentSet map[string]string

func (aSet assignmentSet) toSlice() assignments {
	as := assignments{}
	for fieldName, fieldVal := range aSet {
		as = append(as, assignment{
			fieldName: fieldName,
			fieldVal:  fieldVal,
		})
	}

	sort.Sort(as)
	return as
}

//
// The assignments type.
//

type assignments []assignment

func (a assignments) Len() int {
	return len(a)
}

func (a assignments) Less(i, j int) bool {
	return a[i].fieldName < a[j].fieldName
}

func (a assignments) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//-----------------------------------------------------------------------------
// Property method.
//-----------------------------------------------------------------------------

// propertyMethod represents a function that sets a specific property
// in the Kubernetes API. For example, `v1.Container` has a property,
// `image`, which the user might set in the following way:
//
//   someContainer + container.image("nginx")
//
// This struct represents the definition of all such "property
// methods".
type propertyMethod struct {
	k8s.Property
	name      string
	param     string
	fieldName string
}

func newPropertyMethod(
	name string, param string, fieldName string, property k8s.Property,
) propertyMethod {
	return propertyMethod{
		Property:  property,
		name:      name,
		param:     param,
		fieldName: fieldName,
	}
}

func (pm *propertyMethod) emit(m marshaller) {
	fieldName := pm.fieldName
	if pm.fieldName == errorKeyword {
		fieldName = fmt.Sprintf("\"%s\"", pm.fieldName)
	}

	descLines := strings.Split(pm.Property.Description, "\n")
	for _, line := range descLines {
		m.bufferLine("// " + line)
	}

	// Use a mixin if the constructor set a default value.
	var propMethod string
	if pm.Type != nil && *pm.Type == "array" {
		// Arrays are a special case. If the underlying property is an
		// array, we'll want to check whether the user passed an array or
		// an item. If they passed an array, we want to concatenate the
		// arrays. If they passed an item, we want to append it (which
		// means wrapping it in an array and concatenating that).
		propMethod = fmt.Sprintf(
			"%s(%s):: if std.type(%s) == \"array\" then {%s+: %s} else {%s+: [%s]},",
			pm.name,
			pm.param,
			pm.param,
			fieldName,
			pm.param,
			fieldName,
			pm.param)
	} else {
		var templateString string
		if pm.Ref != nil || (pm.Type != nil && *pm.Type == "object") {
			templateString = "%s(%s):: {%s+: %s},"
		} else {
			templateString = "%s(%s):: {%s: %s},"
		}

		propMethod = fmt.Sprintf(
			templateString,
			pm.name,
			pm.param,
			fieldName,
			pm.param)
	}

	m.bufferLine(propMethod)
}

//
// The methodSet type.
//

type propertyMethodSet map[string]*propertyMethod

func makePropertyMethodSet() propertyMethodSet {
	return make(map[string]*propertyMethod)
}

func (mSet propertyMethodSet) toSlice() propertyMethods {
	ms := propertyMethods{}
	for _, m := range mSet {
		ms = append(ms, m)
	}

	sort.Sort(ms)
	return ms
}

func (mSet propertyMethodSet) addMethod(m *propertyMethod) propertyMethodSet {
	if m != nil {
		mSet[m.name] = m
	}
	return mSet
}

//
// The propertyMethods type.
//

type propertyMethods []*propertyMethod

func (methods propertyMethods) Len() int {
	return len(methods)
}

func (methods propertyMethods) Less(i, j int) bool {
	return methods[i].name < methods[j].name
}

func (methods propertyMethods) Swap(i, j int) {
	methods[i], methods[j] = methods[j], methods[i]
}

//-----------------------------------------------------------------------------
// Mixin.
//-----------------------------------------------------------------------------

// mixinSpec represents the set of mixins associated with a property
// of some `name`. For example, `v1.Service` has a property called
// `metadata`; hence, `mixinSpec` should contain mixins for all the
// properties of the `v1.ObjectMeta` namespace, so that the user can
// do something like:
//
//   someDeployment + deployment.mixin.metadata.name("foo")
//
type mixinSpec struct {
	parent     *namespace
	name       string
	mixinField string
	path       string
	methods    propertyMethodSet
}

func newMixinSpec(
	name, fieldName string, refNS *namespace, parent *namespace,
) *mixinSpec {
	ms := mixinSpec{
		parent:     parent,
		name:       name,
		path:       refNS.path,
		mixinField: fieldName,
		methods:    makePropertyMethodSet(),
	}

	for _, method := range refNS.methods {
		id := method.name
		pm := newPropertyMethod(id, method.param, method.fieldName, method.Property)
		ms.methods[id] = &pm
	}

	return &ms
}

func (ms *mixinSpec) emit(m marshaller) error {
	camelCaseMixinName, err := stylizeString(
		ms.name, ms.parent.ID+".mixin",
		func(newId string) bool {
			_, ok := ms.parent.mixins[newId]
			return ok
		})
	if err != nil {
		return err
	}
	m.bufferLine(fmt.Sprintf("%s:: {", camelCaseMixinName))
	m.indent()

	mixinLocal := makeLocalFunc(
		ms.mixinField,
		"mixin",
		fmt.Sprintf("{%s+: mixin}", ms.mixinField))
	m.bufferLine(mixinLocal)

	for _, mixinMethod := range ms.methods.toSlice() {
		marshalMixinMethod(m, ms.mixinField, "$."+toJsonnetPath(ms.path), mixinMethod)
	}

	m.dedent()
	m.bufferLine("},")
	return nil
}

func marshalMixinMethod(
	m marshaller, mixinFunctionName, mixinPrefix string, method *propertyMethod,
) {
	mixinCall := mixinPrefix + "." + method.name
	var templateString string
	// Use a mixin if the constructor set a default value.
	if method.Ref != nil {
		templateString = "%s(%s):: %s(%s(%s)),"
	} else if method.Type != nil && (*method.Type == "object" || *method.Type == "array") {
		templateString = "%s(%s):: %s(%s(%s)),"
	} else {
		templateString = "%s(%s):: %s(%s(%s)),"
	}

	propMethod := fmt.Sprintf(
		templateString,
		method.name,
		method.param,
		mixinFunctionName,
		mixinCall,
		method.param)
	m.bufferLine(propMethod)
}

//
// The mixinSpecSet type.
//

type mixinSpecSet map[string]*mixinSpec

func makeMixinSpecSet() mixinSpecSet {
	return make(map[string]*mixinSpec)
}

func (msSet mixinSpecSet) toSlice() mixinSpecs {
	mixins := mixinSpecs{}
	for _, mixinSpec := range msSet {
		mixins = append(mixins, mixinSpec)
	}

	sort.Sort(mixins)
	return mixins
}

//
// The mixinSpecs type.
//

type mixinSpecs []*mixinSpec

func (mixins mixinSpecs) Len() int {
	return len(mixins)
}

func (mixins mixinSpecs) Less(i, j int) bool {
	return mixins[i].name < mixins[j].name
}

func (mixins mixinSpecs) Swap(i, j int) {
	mixins[i], mixins[j] = mixins[j], mixins[i]
}
