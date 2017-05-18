# Tutorial

The **ksonnet** [readme][readme] 
shows you how to create a default 
`deployment.json` file that lets you deploy an nginx container 
to an existing Kubernetes cluster. This tutorial shows you how to:

* Modify the deployment using **ksonnet** definitions
* Define other Kubernetes objects
* Work with mixins to develop complex configurations

## Prerequisites

This tutorial assumes that you have performed the following 
tasks. For details, see the [readme][readme].

* Installed **Jsonnet**
* Cloned the **ksonnet** repository locally
* Installed and configured the VisualStudio Code extension 
(optional)
* Created a test Kubernetes cluster

## Modify the default deployment

**ksonnet** lets you configure or modify any Kubernetes object. For 
example, to customize the default `nginx` deployment, you can write:

```c++
local kubeCore = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";
{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json":
    deployment.FromContainer("nginx-deployment", 2, nginxContainer) +
    deployment.mixin.spec.RollingUpdateStrategy() + // add rolling update strategy
    deployment.mixin.spec.Selector({ "app": "nginx" }), // add custom selector
}

Save the file as `customDeploy.libsonnet` 
and run:

```bash
jsonnet customDeploy.libsonnet
kubectl apply -f deployment.json
```

And the generated YAML looks like this. You can see the new `strategy` 
and `selector` fields:

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

## Define other Kubernetes objects

**ksonnet** lets you define any Kubernetes object. For example, 
you can define a container:

```c++
local kubeCore = import "../../kube/core.libsonnet";
local container = kubeCore.v1.container;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80)
```

Save this snippet as container.libsonnet, and run:

```bash
jsonnet container.libsonnet
```

The JSON output looks like this:

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

Or you can include a liveness probe:

```c++
local kubeCore = import "../../kube/core.libsonnet";
local container = kubeCore.v1.container;
local probe = kubeCore.v1.probe;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80) +
container.LivenessProbe(probe.Http("/", 80, 15, 1))
```

Save the file again -- `container.libsonnet` -- 
and run:

```bash
jsonnet container.libsonnet
```

The JSON output now looks like this:

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

Now you can define a pod that runs this container:

```c++
local kubeCore = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

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

Save the file as pod.libsonnet, and run the following commands 
to deploy the pod to your cluster:

```bash
jsonnet pod.libsonnet // create pod.json
kubectl apply -f pod.json
```

## Work with mixins

[readme]: https://github.com/ksonnet/ksonnet-lib/blob/master/README.md "ksonnet readme"


