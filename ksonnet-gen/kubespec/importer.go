package kubespec

import (
	"encoding/json"

	"github.com/go-openapi/spec"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
)

// Import imports an OpenAPI swagger schema.
func Import(path string) (*spec.Swagger, error) {
	b, err := swag.LoadFromFileOrHTTP(path)
	if err != nil {
		return nil, errors.Wrap(err, "load schema from path")
	}

	return CreateAPISpec(b)
}

// CreateAPISpec a swagger file into a *spec.Swagger.
func CreateAPISpec(b []byte) (*spec.Swagger, error) {
	var apiSpec spec.Swagger
	if err := json.Unmarshal(b, &apiSpec); err != nil {
		return nil, errors.Wrap(err, "parse swagger JSON")
	}

	return &apiSpec, nil
}
