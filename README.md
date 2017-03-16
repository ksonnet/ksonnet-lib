# kube.libsonnet

**Concise, correct Kubernetes application definitions, without the YAML.**

`kube.libsonnet`, is a [Jsonnet][jsonnet] library that makes it easy
to build Kubernetes applications by exposing _expressive_,
_flexible_ abstractions for interacting with the [Kubernetes API][v1].

![Jsonnet syntax highlighting][jsonnet-demo]

Most other projects (_e.g._, [Kompose][Kompose],
[OpenCompose][OpenCompose], and [compose2kube][compose2kube]) simplify
the process of creating a Kubernetes application by creating a simpler
API that maps to the Kubernetes API. `kube.libsonnet` instead aims
instead to make it much simpler to build and customize the Kubernetes
API objects themselves, in all their complexity. This results in
concise, modular application definitions, but without losing any of the options
and features of the original Kubernetes API.

For more info, see the following resources:

* **[Hello, stateless world!][hello-world]** A simple example
  application.
* **[Tutorial][tutorial]**, a more in-depth tutorial explaining the
  core abstractions and tools exposed by `kube.libsonnet`.
* **[gitlab.jsonnet][gitlab-jsonnet] and
  [gitlab.libsonnet][gitlab-libsonnet]**. This is a real-world example
  of a Kubernetes application written on `kube.libsonent`. The first
  file is the main entry point for Jsonnet, when it compiles GitLab's
  deployments, services, _etc._, to JSON so that `kubectl` can pick
  them up; the second contains all the nitty-gritty logic of defining
  the each of those components.
* **[Design document][design]**, (_highly incomplete_) explaining the
  goals and rationale behind the core design decisions of
  `kube.libsonnet`.

## Hello, stateless world!

A detailed look at `kube.libsonnet`'s core abstractions is available
in the [tutorial][tutorial]. This section is intended to give just a
taste of what is possible.

Let's start by converting the Kubernetes [nginx hello world
tutorial][helloworld] to use `kube.libsonnet`. Here is the original
YAML:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

And here is the equivalent implementation in Jsonnet with
`kube.libsonnet` (see also [source][v1hellojsonnet]):

```c++
// hello.jsonnet; imports omitted
{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json": deployment.FromContainer("nginx-deployment", 2, nginxContainer),
}
```

Using the Jsonnet command line, we can easily generate a
`deployment.json` file from the above, which can then be sent to the
cluster directly by `kubectl`:

```bash
$ jsonnet hello.jsonnet -m .         # Generates `deployment.json`.
$ kubectl create -f deployment.json  # Dispatch to run on cluster.
```

This is nice, but sometimes the default `Deployment` object is not
precisely what you want.

A core goal of `kube.libsonnet` is to maintain the flexibility and
expressiveness of the original Kubernetes API objects. To give you
some idea of how easy it is to modify these objects, we will show how
we can use `kube.libsonnet`'s mixins to change the deployment
specification to use a rolling update strategy, and a custom selector
(see also [source][v2hellojsonnet]). **The use of mixins is covered in
much more detail during the [tutorial][tutorial];** this is really
intended to give you a taste of what's possible:

```c++
// hello.jsonnet; imports omitted
{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json":
    deployment.FromContainer("nginx-deployment", 2, nginxContainer) +
    deployment.MixinSpec(
      deployment.spec.RollingUpdateStrategy(1, 1) +
      deployment.spec.Selector({ "app": "nginx" })),
}
```

Here we have customized the `strategy` and `selector` fields, but we
can use this method to customize any field in the Kubernetes-standard
[Deployment Spec API object][deploymentspec]. In fact, as you will see
in the [tutorial][tutorial], _all_ `kube.libsonnet` objects are highly
customizable!

Looking at the YAML that is generated, we can begin to understand the
power of these mixins to customize the default API objects:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  strategy:
    rollingUpdate:
        maxSurge: 1,
        maxUnavailable: 1
    type: RollingUpdate
  selector:
    matchLabels:
        app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

## Tutorial

For a detailed tutorial, see [docs/TUTORIAL.md][tutorial]


[jsonnet]: http://jsonnet.org/ "Jsonnet"
[v1]: https://kubernetes.io/docs/api-reference/v1/definitions/ "V1 API objects"
[v1Container]: https://kubernetes.io/docs/api-reference/v1/definitions/#_v1_container "v1.Container"
[Kompose]: https://github.com/kubernetes-incubator/kompose "Kompose"
[OpenCompose]: https://github.com/redhat-developer/opencompose "OpenCompose"
[compose2kube]: https://github.com/kelseyhightower/compose2kube "compose2kube"

[helloworld]: https://kubernetes.io/docs/tutorials/stateless-application/run-stateless-application-deployment/ "Hello, Kubernetes!"
[v1hellojsonnet]: https://github.com/heptio/kube.libsonnet/blob/master/examples/hello-world/hello.v1.jsonnet "Hello, Jsonnet (v1)!"
[v2hellojsonnet]: https://github.com/heptio/kube.libsonnet/blob/master/examples/hello-world/hello.v2.jsonnet "Hello, Jsonnet (v2)!"
[deploymentspec]: https://kubernetes.io/docs/api-reference/extensions/v1beta1/definitions/#_v1beta1_deploymentspec "v1.DeploymentSpec"
[hello-world]: https://github.com/heptio/kube.libsonnet#hello-stateless-world "Hello, stateless world!"
[design]: https://github.com/heptio/kube.libsonnet/blob/master/docs/DESIGN.md "kube.libsonnet design document"
[tutorial]: https://github.com/heptio/kube.libsonnet/blob/master/docs/TUTORIAL.md "kube.libsonnet tutorial"
[gitlab-jsonnet]: https://github.com/heptio/kube.libsonnet/blob/master/examples/kubernetes-gitlab-demo/gitlab-jsonnet/gitlab.jsonnet "gitlab.jsonnet"
[gitlab-libsonnet]: https://github.com/heptio/kube.libsonnet/blob/master/examples/kubernetes-gitlab-demo/gitlab-jsonnet/gitlab.libsonnet "gitlab.libsonent"
[jsonnet-demo]: docs/images/kube-demo.gif
