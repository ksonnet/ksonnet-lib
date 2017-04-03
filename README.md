# kube.libsonnet: concise, correct Kubernetes configurations, without the YAML

By Heptio, Inc., 2017

`kube.libsonnet` provides a simple alternative to writing 
complex YAML for your Kubernetes configurations. It accomplishes
this goal by using the data templating language
[Jsonnet][jsonnet] to write against the
[Kubernetes application API][v1]. This approach also makes it 
easy to extend your configuration as your application scales up.

![Jsonnet syntax highlighting][jsonnet-demo]

Other projects, such as [Kompose][Kompose],
[OpenCompose][OpenCompose], and [compose2kube][compose2kube], simplify
the process of writing a Kubernetes configuration by creating a simpler
API that maps to the Kubernetes API. `kube.libsonnet` instead simplifies 
the work required to build and customize the Kubernetes API objects 
themselves. This approach results in concise, modular configurations, 
without losing any of the options and features of the original 
Kubernetes API.

## Installing and running

First, you need Jsonnet:

`brew install jsonnet`

Then, fork or clone this repository, and add the appropriate import 
statements for the library to your Jsonnet code. For example 
(from the tutorial):

```c++
local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";
```

## Hello, stateless world!

Let's start by converting the Kubernetes 
[nginx hello world tutorial][helloworld] to use `kube.libsonnet`. 
Here is the original YAML:

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
some idea of how easy it is to modify these objects, let's
change the deployment specification to use a rolling update strategy 
and a custom selector (see also [source][v2hellojsonnet]). We use
`kube.libsonnet`'s mixins to make the changes.

```c++
// hello.jsonnet; imports omitted
{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json":
    deployment.FromContainer("nginx-deployment", 2, nginxContainer) +
    deployment.mixin.spec.RollingUpdateStrategy() +
    deployment.mixin.spec.Selector({ "app": "nginx" }),
}
```

Here we customize the `strategy` and `selector` fields, but we
can use this method to customize any field in the Kubernetes-standard
[Deployment Spec API object][deploymentspec]. In fact, as
the [tutorial][tutorial] demonstrates, _all_ `kube.libsonnet` 
objects are customizable.

If you look at the generated YAML, you can see the customized API
objects.

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

## More information

See also the following resources:

* **[Tutorial][tutorial]**. A more in-depth tutorial that explains the
  core abstractions and tools exposed by `kube.libsonnet`.
* **[gitlab.jsonnet][gitlab-jsonnet] and
  [gitlab.libsonnet][gitlab-libsonnet]**. A real-world example
  of a Kubernetes configuration written with `kube.libsonent`. The first
  file is the main entry point for Jsonnet. It compiles GitLab's
  deployments, services, _etc._, to JSON so that `kubectl` can pick
  them up. The second file contains the logic that defines
  each component.
* **[Design document][design]**, (_highly incomplete_) Explains the
  goals and rationale behind the core design decisions.


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
