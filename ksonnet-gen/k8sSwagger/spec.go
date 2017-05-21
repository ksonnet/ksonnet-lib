package k8sSwagger

import (
	"encoding/json"
)

type AppSpec struct {
	SwaggerVersion string           `json:"swaggerVersion"`
	ApiVersion     string           `json:"apiVersion"`
	BasePath       string           `json:"basePath"`
	ResourcePath   string           `json:"resourcePath"`
	Info           Info             `json:"info"`
	Models         map[string]Model `json:"models"`
	// apis
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Model struct {
	ID          string              `json:"id"`
	Description string              `json:"description"`
	Required    []string            `json:"required"`   // Required properties
	Properties  map[string]Property `json:"properties"` // Possible fields
}

type Property struct {
	ObjectType
	Format      *string     `json:"format"` // e.g., int64 if ObjectType is "integer"
	ItemType    *ObjectType `json:"items"`  // if ObjectType.Type == "array"
	Description string      `json:"description"`
}

type ObjectType struct {
	Type *string `json:"type"` // string, array, boolean, object, integer
	Ref  *string `json:"$ref"` // e.g., v1.ObjectMetadata
}

func AppSpecFromJson(swaggerJson []byte) (*AppSpec, error) {
	apiSpec := AppSpec{}
	err := json.Unmarshal(swaggerJson, &apiSpec)
	if err != nil {
		return nil, err
	}

	return &apiSpec, nil
}
