# Backsplice: a simpler way to write concise, correct Kubernetes configurations

By Heptio, Inc., 2017

*Backsplice* provides a simpler alternative to writing 
complex YAML for your Kubernetes configurations. Instead, you 
write template functions against the 
[Kubernetes application API][v1] using the 
data templating language [Jsonnet][jsonnet]
. Components called *mixins* also help
simplify the work that's required to extend your configuration 
as your application scales up.

![Jsonnet syntax highlighting][jsonnet-demo]

Other projects help simplify the work of writing a Kubernetes 
configuration by creating a simpler API that wraps the Kubernetes 
API. These projects include [Kompose][Kompose],
[OpenCompose][OpenCompose], and [compose2kube][compose2kube]. 

*Backsplice* instead streamlines the process of writing 
configurations that create native Kubernetes objects. 

## Install and run

First, install Jsonnet:

`brew install jsonnet`

Then, fork or clone this repository, and add the appropriate import 
statements for the library to your Jsonnet code:

```c++
local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";
```

You might want to consider working in Visual Studio Code, using 
an extension that
provides syntax highlighting and a preview pane for your output
in either YAML or JSON. See <link_to_repo>.

### Get started

If you're not familiar with Jsonnet, check out the website and tutorial. 
The *Backsplice* repository also includes a brief introduction <link>.

You can also start writing `.libsonnet` or `.jsonnet` files based on 
the examples in this readme and in the tutorial <link>. Then run the 
following command:

```bash
jsonnet <filename.libsonnet>
```

## Write your config files with Backsplice

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

Instead, you can write the following *Backsplice*:

```c++
local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";
{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json": deployment.FromContainer("nginx-deployment", 2, nginxContainer),
}
```

Save the file as `helloworld.libsonnet`, then run:

```bash
jsonnet helloword.libsonnet
```

This command creates the `deployment.json` file that the 
*Backsplice* snippet defines.

You can now apply this deployment to your Kubernetes cluster
by running the following command:

```bash
kubectl apply -f deployment.json
```

For more examples and a fuller explanation, see the tutorial <link>

## Contributing

Community contributions are welcome. You can submit a pull request, 
file an issue, or get in touch with us <how? who? where?>

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
