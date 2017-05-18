# ksonnet: a simpler way to write concise, correct Kubernetes configurations

By Heptio, Inc., 2017

**ksonnet** provides a simpler alternative to writing 
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

**ksonnet** instead streamlines the process of writing 
configurations that create native Kubernetes objects. 

## Install and run

First, install Jsonnet:

`brew install jsonnet`

Then, fork or clone this repository, and add the appropriate import 
statements for the library to your Jsonnet code:

```c++
local kubeCore = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";
```

You might want to consider working in Visual Studio Code, using 
an extension that
provides syntax highlighting and a preview pane for your output
in either YAML or JSON. See 
[this GitHub repository](https://github.com/heptio/vscode-jsonnet).

### Get started

If you're not familiar with **Jsonnet**, check out the 
[website](http://jsonnet.org/index.html) and [tutorial]
(http://jsonnet.org/docs/tutorial.html). For usage, see 
the [command line tool]
(http://jsonnet.org/implementation/commandline.html). 
The current repository also includes an 
[introduction to Jsonnet](https://github.com/ksonnet/ksonnet-lib/blob/master/docs/jsonnetIntro.md).

You can also start writing `.libsonnet` or `.jsonnet` files based on 
the examples in this readme and in the [tutorial][tutorial]. Then run the 
following command:

```bash
jsonnet <filename.libsonnet>
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

```c++
local kubeCore = import "../../kube/core.libsonnet";
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
**ksonnet** snippet defines.

You can now apply this deployment to your Kubernetes cluster
by running the following command:

```bash
kubectl apply -f deployment.json
```

## The **ksonnet** libraries

The **ksonnet** libraries provide sets of different methods for 
creating and manipulating Kubernetes objects:

* `kube/core.libsonnet`: extends the object model and functions of `Jsonnet` to implement the Kubernetes API
* `kube/util.libsonnet`: contains methods to help create complex Kubernetes objects out of smaller objects

Kubernetes v1 and v1beta1 are supported.

For more examples and a fuller explanation, see the [tutorial][tutorial].

## Contributing

(content TK from Jorge)

[tutorial]: https://github.com/ksonnet/ksonnet-lib/blob/master/docs/TUTORIAL.md "ksonnet tutorial"

