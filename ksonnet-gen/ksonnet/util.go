package ksonnet

import (
	"errors"
	"fmt"
	"strings"
)

func parsePath(path string) []string {
	return strings.Split(path, ".")
}

func unparsePath(components []string) string {
	return strings.Join(components, ".")
}

func toJsonnetPath(path string) string {
	components := parsePath(path)
	for i, component := range components {
		if len(component) == 0 {
			continue
		}
		camelCaseID := strings.ToLower(string(component[0])) + component[1:]
		components[i] = camelCaseID
	}
	return unparsePath(components)
}

// stylizeID convert an ID is in the style of a Jsonnet ID, i.e., with
// the first letter lowercase. This function will also check that this
// conversion does not trample over another identifier, e.g., that
// `v1.Foo` will not replace an existing identifier, `foo`, whose name
// differs only by case.
func stylizeID(id string, parentNS namespace) (string, error) {
	return stylizeString(id, parentNS.ID, func(newId string) bool {
		_, ok := parentNS.namespaceSet[newId]
		return ok
	})
}

func stylizeString(
	s string, nsID string, checkCollision func(string) bool,
) (string, error) {
	if len(s) < 1 {
		msg := fmt.Sprintf(
			"Object in namespace '%s' has an ID of 0 characters", nsID)
		return "", errors.New(msg)
	}

	camelCaseID := strings.ToLower(string(s[0])) + s[1:]
	// NOTE: Check if `currNS.ID` != `camelCaseID` to guard against
	// case when the first letter was lowercase to begin with.
	if s != camelCaseID && checkCollision(camelCaseID) {
		msg := fmt.Sprintf(
			"Can't lower-case first letter of '%s' in namespace '%s' because a member already exists with that name",
			s,
			nsID)
		return "", errors.New(msg)
	}

	return camelCaseID, nil
}

func makeLocalObjectVal(id, field, val string) string {
	return fmt.Sprintf("local %s = {%s: \"%s\"},", id, field, val)
}

func makeLocalFunc(id, params, body string) string {
	return fmt.Sprintf("local %s(%s) = %s,", id, params, body)
}
