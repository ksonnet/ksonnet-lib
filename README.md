# ksonnet: a simpler way to write concise, correct Kubernetes configurations

By Heptio, Inc., 2017

*ksonnet* provides a simpler alternative to writing 
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

*ksonnet* instead streamlines the process of writing 
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
The *ksonnet* repository also includes a brief introduction <link>.

You can also start writing `.libsonnet` or `.jsonnet` files based on 
the examples in this readme and in the tutorial <link>. Then run the 
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

Instead, you can write the following *ksonnet* code:

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
*ksonnet* snippet defines.

You can now apply this deployment to your Kubernetes cluster
by running the following command:

```bash
kubectl apply -f deployment.json
```

## The *ksonnet* files

*ksonnet* comprises two files that contain different methods for 
creating and manipulating Kubernetes objects:

* `kubeCore.libsonnet`: extends the object model and functions of `Jsonnet` to implement the Kubernetes API
* `kubeUtil.libsonnet`: contains methods to help create complex Kubernetes objects out of smaller objects

Kubernetes v1 and v1beta1 are supported.

For more examples and a fuller explanation, see the tutorial <link>

