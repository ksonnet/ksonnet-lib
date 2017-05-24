# ksonnet: a simpler way to write concise, correct Kubernetes configurations

By Heptio, Inc., 2017

**ksonnet** (currently in beta testing) provides a simpler alternative
to writing complex YAML for your Kubernetes configurations. Instead,
you write template functions against the [Kubernetes application
API][v1] using the data templating language [Jsonnet][jsonnet] .
Components called **mixins** also help simplify the work that's
required to extend your configuration as your application scales up.

![Jsonnet syntax highlighting][jsonnet-demo]

Other projects help simplify the work of writing a Kubernetes
configuration by creating a simpler API that wraps the Kubernetes
API. These projects include [Kompose][Kompose],
[OpenCompose][OpenCompose], and [compose2kube][compose2kube].

**ksonnet** instead streamlines the process of writing
configurations that create native Kubernetes objects.

## Install

First, install Jsonnet.

### Mac OS X

If you do not have Homebrew installed, [install it now](https://brew.sh/).

Then run:

`brew install jsonnet`

### Linux

You must build the binary. For details, [see the GitHub
repository](https://github.com/google/jsonnet).

## Run

Fork or clone this repository, using a command such as:

```shell
git clone git@github.com:ksonnet/ksonnet-lib.git
```

Then add the appropriate import
statements for the library to your Jsonnet code:

```javascript
local k = import "ksonnet.beta.1/k.libsonnet";
```

As Jsonnet [does not yetsupport](https://github.com/google/jsonnet/issues/9)
`go get`-style importing from HTTP, import paths are relative to the root of the
**ksonnet** repository. Remember to modify the paths appropriately
when you work in another environment so that they point at your clone.

Additionally, since Jsonnet does not yet support [ES2016-style](https://github.com/google/jsonnet/issues/307) imports it is common to "unpack" an import with a series of `local` definitions, as so:

```javascript
local container = k.core.v1.container;
local deployment = k.extensions.v1beta1.deployment;
```

### Tools

Developed in tandem with `ksonnet-lib` is
[`vscode-jsonnet`](https://github.com/heptio/vscode-jsonnet), a static
analysis toolset written as a [Visual Studio
Code](https://code.visualstudio.com/) plugin, meant to provide
features such as autocomplete, syntax highlighting, and static
analysis.

### Get started

If you're not familiar with **Jsonnet**, check out the
[website](http://jsonnet.org/index.html) and
[their tutorial](http://jsonnet.org/docs/tutorial.html). For usage, see
the [command line tool](http://jsonnet.org/implementation/commandline.html).
This repository also includes an
[introduction to Jsonnet](docs/jsonnetIntro.md).

You can also start writing `.libsonnet` or `.jsonnet` files based on
the examples in this readme and in the [tutorial][tutorial]. Then run the
following command:

```bash
jsonnet <filename.libsonnet>
```

This command produces a JSON file that you can then run the
appropriate `kubectl`
commands against, with the following syntax:

```bash
kubectl <command> -<options> <filename.json>
```

## Write your config files with ksonnet

The YAML for the Kubernetes
[nginx hello world tutorial][helloworld] looks
like this:

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

Instead, you can write the following **ksonnet** code:

```javascript
local k = import "../../ksonnet.beta.1/k.libsonnet";

local container = k.core.v1.container;
local deployment = k.apps.v1beta1.deployment;

local nginxContainer =
  container.default("nginx", "nginx:1.7.9") +
  container.helpers.namedPort("http", 80);

deployment.default("nginx-deployment", nginxContainer) +
deployment.mixin.spec.replicas(2)
```

Save the file as `helloworld.libsonnet`, then run:

```bash
jsonnet helloworld.libsonnet
```

This command creates the `deployment.json` file that the
**ksonnet** snippet defines.

You can now apply this deployment to your Kubernetes cluster
by running the following command:

```bash
kubectl apply -f deployment.json
```

For a full-scale example, compare the [**ksonnet** definition for
the CockroachDB Helm chart][cockroachks] with the
[original YAML][cockroachch].

## The **ksonnet** libraries

The **ksonnet** project organizes libraries by the level of
abstraction they approach. For most users, the right entry point is:

* `ksonnet.beta.1/k.libsonnet`: higher-level abstractions and methods
  to help create complex Kubernetes objects out of smaller objects

`k.libsonnet` is built on top of several lower-level utility libraries
that extends the object model and functions of `Jsonnet` to implement
the Kubernetes API. These are generated directly from the swagger
spec, and are organized by the API groups they belong to:

* `ksonnet.beta.1/apps.v1beta1.libsonnet`
* `ksonnet.beta.1/core.v1.libsonnet`
* `ksonnet.beta.1/extensions.v1beta1.libsonnet`

For more examples and a fuller explanation, see the [tutorial][tutorial].

## Contributing

Thanks for taking the time to join our community and start
contributing!

### Before you start

* Please familiarize yourself with the [Code of
Conduct](CODE-OF-CONDUCT.md) before contributing.
* See [CONTRIBUTING.md](CONTRIBUTING.md) for instructions on the
developer certificate of origin that we require.

### Pull requests

* We welcome pull requests. Feel free to dig through the
[issues](https://github.com/ksonnet/ksonnet-lib/issues) and jump in.

## Contact us

Have any questions or long-form feedback? You can always find us here:

* Our [Slack channel](https://ksonnet.slack.com).
* Our [mailing list](https://groups.google.com/forum/#!forum/ksonnet).
* We monitor the [ksonnet
tag](https://stackoverflow.com/questions/tagged/ksonnet) on Stack
Overflow.

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

[jsonnet-demo]: docs/images/kube-demo.gif

[tutorial]: docs/TUTORIAL.md "ksonnet tutorial"
[cockroachks]: examples/charts/cockroachdb/cockroachdb.jsonnet "cockroachdb ksonnet"
[cockroachch]: https://github.com/kubernetes/charts/blob/master/stable/cockroachdb/templates/cockroachdb-petset.yaml "cockroachdb YAML"
