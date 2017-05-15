# Tutorial



In this document, we will take a guided tour through the
`kube.libsonnet` library, covering the core abstractions exposed by
the library, and providing some examples of how to use Jsonnet's
more powerful features to create extensible Kubernetes applications.

For more background on Jsonnet, see excellent [official
tutorial][jsonnetTutorial].

|    | Section                                                   |
|----|-----------------------------------------------------------|
| 1  | [Jsonnet recap][jsonnet-recap]                            |
| 2  | [kubeCore.libsonnet][kubeCore]                            |
| 2a | [Improved support for Kubernetes primitives][k8s-prims]   |
| 2b | [Simple, static schema validation][schema-validation]     |
| 2c | [Flexible customization via mixins][customization-mixins] |
| 3  | [kubeUtil.libsonnet][kubeUtil]                            |
| 3a | [Gitlab: A real-world example][gitlab-example]            |


# 2: `kubeCore.libsonnet`

The `kube.libsonnet` project (temporarily) does not currently have a
single `kube.libsonnet` file. Instead, the project is currently split
into two files:

* `kubeCore.libsonnet`, which uses Jsonnet's powerful object model and
  function support to (concisely) implement the Kubernetes API (both
  [v1][v1] and [v1beta1][v1beta1]); and
* `kubeUtil.libsonnet`, which is a small collection of helper methods
  intended to make it easier to glue Kubernetes API objects together.

This section covers `kubeCore.libsonnet`. We will incidentally use
`kubeUtil.libsonnet` to aid in the discussion, but the vast majority
of the ins and outs will be covered in the [kubeUtil
section][kubeUtil].

## 2a: Improved support for Kubernetes primitives

`kubeCore.libsonnet` exposes a collection of abstractions whose goal,
in aggregate, is to make it vastly easier to incrementally build and
directly manipulate the core Kubernetes API objects.

In the following example, we use the `kubeCore` API to up a simple
[v1.Container][v1-container] object, which we can then add to a Pod
definition:

```c++
// examples/tutorial/simpleContainer.1.jsonnet
local kubeCore = import "../../kubeCore.libsonnet";
local container = kubeCore.v1.container;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80),
```

Observing the output of `jsonnet` below, we can see that this code
produces a normal, plain-old `v1.Container` Kubernetes API object, one
that we can actually include in a pod definition:

```json
{
   "image": "nginx:1.7.9",
   "imagePullPolicy": "Always",
   "name": "nginx",
   "ports": [
      {
         "containerPort": 80,
         "name": "http"
      }
   ]
}
```

`kubeCore` implements all of the `v1.Container` API, so to customize
other properties of the container, you need only use the `+` operator
to add additional properties.

In the following example, we use the `+` operator to add a liveness
probe:

```c++
// examples/tutorial/simpleContainer.2.jsonnet
local kubeCore = import "../../kubeCore.libsonnet";
local container = kubeCore.v1.container;
local probe = kubeCore.v1.probe;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80) +
container.LivenessProbe(probe.Http("/", 80, 15, 1))
```

`jsonnet` then generates the following:

```json
{
   "image": "nginx:1.7.9",
   "imagePullPolicy": "Always",
   "livenessProbe": {
      "httpGet": {
         "path": "/",
         "port": 80
      },
      "initialDelaySeconds": 15,
      "timeoutSeconds": 1
   },
   "name": "nginx",
   "ports": [
      {
         "containerPort": 80,
         "name": "http"
      }
   ]
}
```

Similar such constructs exist for [NOTE: work in progress, this is
kind of a lie!] the entire Kubernetes API.

Since this is a normal, `v1.Container`, we can additionally pass this
to a pod definition. For sake of illustration, here is an example of
how we might do this using a utility function in `kubeUtil` that
creates a pod definition from a `v1.Container` object:

```c++
// examples/tutorial/simplePod.jsonnet
local kubeCore = import "../../kubeCore.libsonnet";
local kubeUtil = import "../../kubeUtil.libsonnet";

local container = kubeCore.v1.container;
local probe = kubeCore.v1.probe;
local pod = kubeUtil.app.v1.pod;

{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80) +
    container.LivenessProbe(probe.Http("/", 80, 15, 1)),

  "nginxPod.json": pod.FromContainer(nginxContainer),
}
```

Running something like `jsonnet examples/tutorial/simplePod.jsonnet -m
.` will produce an `nginxPod.json` containing the JSON below; notice
that this JSON is a normal, well-formed [v1.Pod][v1-pod] object, which
uses the container object from above, and which can be passed to
`kubectl` to run on an arbitrary kubernetes cluster.

```json
{
    "apiVersion": "v1",
    "kind": "Pod",
    "meadata": {
        "labels": {
          "app": "nginx"
        }
    },
    "spec": {
        "containers": [
          {
              "image": "nginx:1.7.9",
              "imagePullPolicy": "Always",
              "livenessProbe": {
                "httpGet": {
                    "path": "/",
                    "port": 80
                },
                "initialDelaySeconds": 15,
                "timeoutSeconds": 1
              },
              "name": "nginx",
              "ports": [
                {
                    "containerPort": 80,
                    "name": "http"
                }
              ]
          }
        ]
    }
}
```

## 2b: Simple, static schema validation

Correct-by-default, etc.

## 2c: Flexible customization via mixins

# 3: `kubeUtil.libsonnet`

[intro here]

## 3a: Gitlab: A real-world example


[v1]: https://kubernetes.io/docs/api-reference/v1/definitions/ "Kubernetes v1 API"
[v1beta1]: https://kubernetes.io/docs/api-reference/extensions/v1beta1/definitions/ "Kubernetes v1beta1 API"
[v1-container]: https://kubernetes.io/docs/api-reference/v1/definitions/#_v1_container "v1.Container"
[v1-pod]: https://kubernetes.io/docs/api-reference/v1/definitions/#_v1_pod "v1.Pod"

[jsonnetTutorial]: http://jsonnet.org/docs/tutorial.html "Jsonnet tutorial"
[jsonnetSyntax]: http://jsonnet.org/docs/tutorial.html#syntax_improvements "Jsonnet syntax improvements"
[jsonnetFunctions]: http://jsonnet.org/docs/tutorial.html#functions "Jsonnet functions"
[jsonnetLocals]: http://jsonnet.org/docs/tutorial.html#locals "Jsonnet local variables"
[jsonnetReferences]: http://jsonnet.org/docs/tutorial.html#references "Jsonnet references"
[jsonnetImports]: http://jsonnet.org/docs/tutorial.html#imports "Jsonnet imports"
[jsonnetOO]: http://jsonnet.org/docs/tutorial.html#oo "Jsonnet OO"
[jsonnetAlgebra]: http://jsonnet.org/language/spec.html#properties "Jsonnet inheritance algebra"

[jsonnet-recap]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#1-jsonnet-recap-references-variables-simple-json-templating "Jsonnet recap"
[kubeCore]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#2-kubecorelibsonnet "kubeCore.libsonnet"
[k8s-prims]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#2a-improved-support-for-kubernetes-primitives "Kubernetes Primitives"
[schema-validation]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#2b-simple-static-schema-validation "Schema validation"
[customization-mixins]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#2c-flexible-customization-via-mixins "Customization with mixins"

[kubeUtil]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#3-kubeutillibsonnet "kubeUtil.libsonnet"
[gitlab-example]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md#3a-gitlab-a-real-world-example "Gitlab: A real-world example"
