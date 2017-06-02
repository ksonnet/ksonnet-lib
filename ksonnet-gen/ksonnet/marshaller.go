package ksonnet

import (
	"bytes"
	"fmt"
	"strings"
)

// marshaller abstracts the task of writing out indented text to a
// buffer. Different components can call `indent` and `dedent` as
// appropriate to specify how indentation needs to change, rather than
// to keep track of the current indentation.
//
// For example, if one component is responsible for writing an array,
// and an element in that array is a function, the component
// responsible for the array need only know to call `indent` after the
// '[' character and `dedent` before the ']' character, while the
// routine responsible for writing out the function can handle its own
// indentation independently.
type marshaller struct {
	depth  int
	prefix string
	lines  []string
	buffer *bytes.Buffer
}

func newMarshaller() *marshaller {
	var buffer bytes.Buffer
	return &marshaller{
		depth:  0,
		prefix: "",
		lines:  []string{},
		buffer: &buffer,
	}
}

func (m *marshaller) bufferLine(text string) {
	line := fmt.Sprintf("%s%s\n", m.prefix, text)
	m.lines = append(m.lines, line)
}

func (m *marshaller) writeAll() ([]byte, error) {
	for _, line := range m.lines {
		_, err := m.buffer.WriteString(line)
		if err != nil {
			return nil, err
		}
	}

	return m.buffer.Bytes(), nil
}

func (m *marshaller) indent() {
	m.depth++
	m.prefix = strings.Repeat("  ", m.depth)
}

func (m *marshaller) dedent() {
	m.depth--
	m.prefix = strings.Repeat("  ", m.depth)
}
