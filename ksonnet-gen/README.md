# ksonnet-gen

`ksonnet-gen` takes the OpenAPI Kubernetes specification and generates
a Jsonnet file representing that API definition.

## Build

```bash
dep ensure      # Fetch dependencies
go build -o ksonnet-gen .
```

## Usage

`ksonnet-gen [path to k8s OpenAPI swagger.json] [output dir]`

Typically the swagger spec is in something like
`k8s.io/kubernetes/api/openapi-spec`, where `k8s.io` is in your Go src
folder.
