# Tutorial

In this document, we will take a guided tour through the
`ksonnet-lib` library, covering the core abstractions exposed by
the library, and providing some examples of how to use Jsonnet's
more powerful features to create extensible Kubernetes applications.

For more background on Jsonnet, see excellent [official
tutorial][jsonnetTutorial].

|    | Section                                                   |
|----|-----------------------------------------------------------|
| 1  | [Jsonnet recap][jsonnet-recap]                            |
| 2  | [core.libsonnet][ksonnet-core]                            |
| 2a | [Improved support for Kubernetes primitives][k8s-prims]   |
| 2b | [Simple, static schema validation][schema-validation]     |
| 2c | [Flexible customization via mixins][customization-mixins] |
| 3  | [kubeUtil.libsonnet][ksonnet-util]                            |
| 3a | [Gitlab: A real-world example][gitlab-example]            |

# 1: Jsonnet recap: references, variables, simple JSON templating

Before we demonstrate the core abstractions of `ksonnet-lib`, it is
worth spending time on a whirlwind tour of Jsonnet, to familiarize
ourselves with the features we will use in this tutorial. (_If you
already know Jsonnet, you will lose nothing by skipping to the [next
section][k8s-prims]._)

For the purposes of this tutorial, you can think of Jsonnet as a
domain-specific language meant to make it easy to declare and template
languages. Think JSON, but with:

* variables (both [lexically-scoped locals][jsonnetLocals] and
  JsonPath-style [references][jsonnetReferences])
* [functions][jsonnetFunctions]
* the ability to define libraries and [import][jsonnetImports] them
* some notion of [object-oriented inheritance between JSON
  objects][jsonnetOO]
* and a bunch of the [syntax annoyances ironed out][jsonnetSyntax].

For the purposes of this tutorial, you need only a very small subset
of these concepts. They are:

### Local variables and references

In Jsonnet, it is possible to define lexically-scoped local variables:

```c++
{
  local foo = "bar",
  baz: foo,
}
```

which produces:

```json
{ "baz": "bar" }
```

Jsonnet additionally exposes a `self` to access properties of the
current object, and a JsonPath-style `$`, which refers to the "root
object" (or: the grandparent who is farthest away from the `$`):

```c++
{
  foo: "bar",
  baz: self.foo,
  cow: {
    moo: $.foo,
  },
}
```

```json
{
  "foo": "bar",
  "baz": "bar",
  "cow": { "moo": "bar" }
}
```

It is worth noting that both `local` variables and references are
_order-independent_, which is a decision that largely falls out of
JSON's design. Notice, for example, that if we re-order `foo` and
`baz`, it does not affect the output of Jsonnet:

```c++
{
  baz: self.foo,
  cow: {
    moo: $.foo,
  },

  // This is perfectly legal.
  foo: "bar",
}
```

### Functions

Jsonnet implements lexically-scoped functions, but they can be
declared in a few ways, and it's worth pointing them out.

In the example below, note the use of two semicolons (_i.e._, `::`) in
the declaration of `function2`. This marks the field as _hidden_,
which is a concept we will look closer at in the section on
object-orientation. For now, it is only important to understand that a
function must be either `local` or hidden with `::`, because Jsonnet
doesn't know how to render a function as JSON data. (Instead of
rendering it, Jsonnet will complain and crash.)

```c++
{
  local function1(arg1) = { foo: arg1 },
  function2(arg1="cluck"):: { bar: arg1 },
  cow: function1("moo"),
  chicken: self.function2(),
}
```

```json
{
   "chicken": {
      "bar": "cluck"
   },
   "cow": {
      "foo": "moo"
   }
}
```

### Object-orientation (inheritance, mixins, _etc_.)

One of Jsonnet's most powerful features, which we use liberally in
this tutorial and in `ksonnet-lib`, is its object model, which
implements a concise, [well-specified _algebra_][jsonnetAlgebra] for
combining JSON-like objects.

The primary tool for combining objects is the `+` operator. In this
example we see two objects (the first is called the _parent_, or
_base_, and the second is called the _child_) that are combined with
the `+`. The child (which is said to _inherit_ from the parent)
overwrites the `bar` property that was defined in the parent:

```c++

{
  // Parent object.
  foo: "foo",
  bar: "bar",
} + {
  // Child object.
  bar: "fubar",
}
```

```json
{
   "bar": "fubar",
   "foo": "foo"
}
```

It is sometimes convenient for a child to reference members of the
parent, so Jsonnet also exposes `super`, which behaves a lot like
`self`, except in reference to the parent:

```c++
{
  foo: "foo",
} + {
  bar: super.foo + "bar",
}
```

```json
{
   "bar": "foobar",
   "foo": "foo"
}
```

One interesting aspect of `super` is that it can be "mixed in",
meaning that if you have an object that refers to `super.bar`, then it
can dynamically be made to inherit from _any object_ that has a `bar`
property. For example:

```c++
local fooTheBar = { bar: super.bar + "foo" };
{
  bar: "bar",
} + fooTheBar
```

```json
{
   "bar": "barfoo"
}
```

This stands in contrast to the object model of (say) Java, where you
would have to declare at compile time an `Animal` class before a `Dog`
class could be made to inherit from it. The technique above (called a
_mixin_) causes the object to inherit dynamically, at runtime rather
than compile time.

Lastly, Jsonnet allows you to create hidden properties, not included
when we generate the final JSON. Denoted with with a `::`, they are
also visible to all descendent objects (_i.e._, children,
grandchildren, _etc_.), and are useful for holding data you'd like to
use to construct other properties, but not expose as part of the
generated JSON itself:

```c++
{
  foo:: "foo",
} + {
  bar: super.foo + "bar",
}
```

```json
{
   "bar": "foobar"
}
```

# 2: `core.libsonnet`

The `ksonnet-lib` project (temporarily) does not currently have a
single `libsonnet` file. Instead, the project is currently split
into two files:

* `core.libsonnet`, which uses Jsonnet's powerful object model and
  function support to (concisely) implement the Kubernetes API (both
  [v1][v1] and [v1beta1][v1beta1]); and
* `kubutileUtil.libsonnet`, which is a small collection of helper methods
  intended to make it easier to glue Kubernetes API objects together.

This section covers `core.libsonnet`. We will incidentally use
`util.libsonnet` to aid in the discussion, but the vast majority
of the ins and outs will be covered in the [kubeUtil
section][ksonnet-util].

## 2a: Improved support for Kubernetes primitives

`core.libsonnet` exposes a collection of abstractions whose goal,
in aggregate, is to make it vastly easier to incrementally build and
directly manipulate the core Kubernetes API objects.

In the following example, we use the `core` API to up a simple
[v1.Container][v1-container] object, which we can then add to a Pod
definition:

```c++
// examples/tutorial/simpleContainer.1.jsonnet
local core = import "../../core.libsonnet";
local container = core.v1.container;

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

`core` implements all of the `v1.Container` API, so to customize
other properties of the container, you need only use the `+` operator
to add additional properties.

In the following example, we use the `+` operator to add a liveness
probe:

```c++
// examples/tutorial/simpleContainer.2.jsonnet
local core = import "../../core.libsonnet";
local container = core.v1.container;
local probe = core.v1.probe;

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
how we might do this using a utility function in `util` that
creates a pod definition from a `v1.Container` object:

```c++
// examples/tutorial/simplePod.jsonnet
local core = import "../../core.libsonnet";
local util = import "../../util.libsonnet";

local container = core.v1.container;
local probe = core.v1.probe;
local pod = util.app.v1.pod;

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

# 3: `util.libsonnet`

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

[jsonnet-recap]: #1-jsonnet-recap-references-variables-simple-json-templating "Jsonnet recap"
[ksonnet-core]: #2-corelibsonnet "ksonnet Core"
[k8s-prims]: #2a-improved-support-for-kubernetes-primitives "Kubernetes Primitives"
[schema-validation]: #2b-simple-static-schema-validation "Schema validation"
[customization-mixins]: #2c-flexible-customization-via-mixins "Customization with mixins"

[kubeUtil]: #3-utillibsonnet "util.libsonnet"
[gitlab-example]: #3a-gitlab-a-real-world-example "Gitlab: A real-world example"
