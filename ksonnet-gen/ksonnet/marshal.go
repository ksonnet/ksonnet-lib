package ksonnet

import (
	"bytes"
	"fmt"
	"strings"

	k8s "github.com/ksonnet/ksonnet-lib/ksonnet-gen/k8sSwagger"
)

// Marshal takes an `AppSpec` (generated from, say, `v1.json` or
// `apps_v1beta1.json`), and creates a Jsonnet file for use as a
// component of `ksonnet-lib`.
func Marshal(appSpec *k8s.AppSpec) ([]byte, error) {
	b := newBuilder(appSpec.ApiVersion)
	b.AddAppSpec(appSpec)

	currNS := b.Root()

	m := newCoreMarshaller()
	if err := m.marshal(nil, currNS); err != nil {
		return nil, err
	}

	bytes, err := m.writeAll()
	if err == nil {
		return bytes, nil
	}
	return nil, err
}

type marshaller interface {
	bufferLine(string)
	writeAll() ([]byte, error)
	indent()
	dedent()
	marshal(parent *namespace, current *namespace) error
}

// coreMarshaller turns a tree of namespaces into a `core.libsonnet`
// file.
type coreMarshaller struct {
	depth  int
	prefix string
	lines  []string
	buffer *bytes.Buffer
}

func newCoreMarshaller() marshaller {
	var buffer bytes.Buffer
	return &coreMarshaller{
		depth:  0,
		prefix: "",
		lines:  []string{},
		buffer: &buffer,
	}
}

func (m *coreMarshaller) bufferLine(text string) {
	line := fmt.Sprintf("%s%s\n", m.prefix, text)
	m.lines = append(m.lines, line)
}

func (m *coreMarshaller) writeAll() ([]byte, error) {
	for _, line := range m.lines {
		_, err := m.buffer.WriteString(line)
		if err != nil {
			return nil, err
		}
	}

	return m.buffer.Bytes(), nil
}

func (m *coreMarshaller) indent() {
	m.depth++
	m.prefix = strings.Repeat("  ", m.depth)
}

func (m *coreMarshaller) dedent() {
	m.depth--
	m.prefix = strings.Repeat("  ", m.depth)
}

// marshalHelper is a recursive helper for `Marshal`. NOTE: `parentNS`
// is `nil` if this is the root namespace
func (m *coreMarshaller) marshal(
	parentNS *namespace, currNS *namespace,
) error {
	//
	// First, emit the object declaration; for an object called
	// `v1.Foo`, this will usually look like `foo:: {`. The concersion
	// of the first letter to lower case is so that the generated code
	// will fit Jsonnet style.
	//

	if parentNS == nil {
		// Root namespace doesn't have a name, so the declaration should
		// just be the string `{` (as opposed to something like
		// `foo:: {`).
		m.bufferLine("{")
	} else {
		camelCaseID, err := stylizeID(currNS.ID, *parentNS)
		if err != nil {
			return err
		}

		m.bufferLine(fmt.Sprintf("%s:: {", camelCaseID))
	}
	m.indent()

	//
	// Second, emit:
	// * Optional `local` that holds the `kind` of the object (e.g.,
	//   `v1.Service` will use the string "Service" as its kind, while
	//   many objects do not have kind field).
	// * Optional `local` that holds the API version for the object
	//   (this should only exist in objects like `v1:: {`, which
	//   represent the v1 namespace).
	// * Definitions of all methods
	//

	if _, ok := currNS.specialMethods[objectKind]; ok {
		idPath := strings.Split(currNS.ID, ".")
		if len(idPath) != 1 {
			return fmt.Errorf(
				"Tried to find `kind` of object, but namespace ID seems to be malformatted: '%s'",
				currNS.ID)
		}
		m.bufferLine(makeLocalObjectVal(objectKind, objectKind, idPath[0]))
	}

	// The root module.
	if currNS.APIVersion != nil {
		m.bufferLine(makeLocalObjectVal(
			objectAPIVersion, objectAPIVersion, *currNS.APIVersion))

		params := strings.Join([]string{metadataName, metadataNamespace}, ", ")
		metadataNameCall := fmt.Sprintf(
			"$.%s.%s(%s)", objectMetadataID, metadataName, metadataName)
		metadataNamespaceCall := fmt.Sprintf(
			"$.%s.%s(%s)", objectMetadataID, metadataNamespace, metadataNamespace)
		body := fmt.Sprintf(
			"{%s: %s + %s}", objectMetadata, metadataNameCall, metadataNamespaceCall)
		m.bufferLine(makeLocalFunc(defaultMetadataFunctionName, params, body))
	}

	if currNS.constructor != nil {
		err := currNS.constructor.emit(m)
		if err != nil {
			return err
		}
	}

	for _, method := range currNS.methods.toSlice() {
		method.emit(m)
	}

	//
	// Third, generate mixins.
	//

	m.marshalMixins(currNS)

	//
	// Third, recursively emit all child namespaces.
	//

	for _, newNS := range currNS.namespaceSet.toSlice() {
		err := m.marshal(currNS, newNS)
		if err != nil {
			return err
		}
	}

	//
	// Fourth, dedent, emit closing brace, return.
	//

	m.dedent()

	if parentNS == nil {
		// Don't emit trailing comma if it's root object.
		m.bufferLine("}")
	} else {
		m.bufferLine("},")
	}

	return nil
}

func (m *coreMarshaller) marshalMixins(currNS *namespace) error {
	// Check whether the `mixin` body would be empty. Because of the
	// small number of mixins we expect the cost of scanning mixins
	// twice to be negligible, and not worth the effort of rearranging
	// the emission code below to first collect and then emit all
	// relevant lines iff the mixin body would not be empty.
	allMixinsEmpty := true
	for _, mixinSpec := range currNS.mixins {
		if len(mixinSpec.methods) != 0 {
			allMixinsEmpty = false
		}
	}
	if allMixinsEmpty {
		return nil
	}

	m.bufferLine("mixin:: {")
	m.indent()

	for _, mixinSpec := range currNS.mixins.toSlice() {
		if len(mixinSpec.methods) == 0 {
			continue
		}

		mixinSpec.emit(m)

	}

	m.dedent()
	m.bufferLine("},")

	return nil
}
