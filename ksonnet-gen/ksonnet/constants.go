package ksonnet

// constructorName is the name `ksonnet-lib` gives to constructors.
var constructorName = "default"

// errorKeyword represents the keyword `error` in Jsonnet. Useful for
// checking whether a property method collides with this keyword.
var errorKeyword = "error"

// objectMetadataID is the path in the Kubernetes API spec given to
// metadata.
// TODO: Verify that this exists before deserializing.
var objectMetadataID = "v1.objectMeta"

// defaultMetadataFunctionName is the name `ksonnet-lib` gives to the
// `local` helper function that constructs a default metadata object.
// This is useful because the spec says that every object with an
// `apiVersion` should have a `name` and a `namespace`, but these are
// not required fields in `v1.objectMetadata`.
var defaultMetadataFunctionName = "defaultMetadata"

// objectAPIVersion is the string the Kubernetes API spec uses in
// objects to denote the version of the API.
var objectAPIVersion = "apiVersion"

// kind is the string the Kubernetes API spec uses in objects to
// denote the kind of object it is.
var objectKind = "kind"

// metadata is the string the Kubernetes API spec uses in objects to
// define the object metadata.
var objectMetadata = "metadata"

// name is the string the Kubernetes API spec uses in objects to
// define the `name` field in object metadata.
var metadataName = "name"

// namespaec is the string the Kubernetes API spec uses in objects to
// define the `namespace` field in object metadata.
var metadataNamespace = "namespace"

// specialResourcePropertiesSet is a set of strings that denote
// "special" properties. For example, `"apiVersion"`.
var specialResourcePropertiesSet = map[string]bool{
	objectAPIVersion: true,
	objectKind:       true,
	objectMetadata:   true,
}
