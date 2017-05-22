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

NOTE: All import paths are relative to the root of the 
*ksonnet** repository.

## Modify the default deployment

**ksonnet** lets you configure or modify any Kubernetes object. For 
example, to customize the default `nginx` deployment, you can write:

```javascript
local core = import "../../kube/core.libsonnet";
local util = import "../../kube/util.libsonnet";

local container = core.v1.container;
local deployment = util.app.v1beta1.deployment;

{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json":
    deployment.FromContainer("nginx-deployment", 2, nginxContainer) +
    deployment.mixin.spec.RollingUpdateStrategy() + // add rolling update strategy
    deployment.mixin.spec.Selector({ "app": "nginx" }), // add custom selector
}
```

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

```javascript
local core = import "../../kube/core.libsonnet";
local container = core.v1.container;

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

```javascript
local core = import "../../kube/core.libsonnet";
local container = core.v1.container;
local probe = core.v1.probe;

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

```javascript
local core = import "../../kube/core.libsonnet";
local util = import "../../kube/util.libsonnet";

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

Save the file as pod.libsonnet, and run the following commands 
to deploy the pod to your cluster:

```bash
jsonnet pod.libsonnet // create pod.json
kubectl apply -f pod.json // apply pod definition to cluster
```

## Work with mixins

You've seen how to modify the default deployment by writing 
**mixins** to add custom `strategy` and `selector` fields. 
As the Jsonnet tutorial explains in more detail, mixins provide 
dynamic inheritance, at runtime instead of compile time. This 
approach means that different team members can define the 
Kubernetes objects that they need. You can then mix them into 
your master definition without having to copy all the details.

For example, you could write the following code to define a 
container for your application:

```javascript
local core = import "../../kube/core.libsonnet";
local util = import "../../kube/util.libsonnet";

local container = core.v1.container;

{
  local appContainer =
       container.Default(app.name, config.containerImage) +
       container.Command(command) +
       container.Ports([
           port.container.Default(8102),
           port.container.Default(9102),
       ]) +
       container.LivenessProbe(probe.Http("/ping", 9102, 10, 2)) +
       container.ReadinessProbe(probe.Http("/ready-to-serve", 9102, 10, 2)) +
       container.Env([
           environ.ValueFromFieldRef("POD_NAME", "metadata.name"),
           environ.ValueFromFieldRef("POD_NAMESPACE", "metadata.namespace"),
           environ.Variable("SERVICE_NAME", app.name),
           environ.Variable("DUMMY_VAR_FOR_NO_OP_DEPLOYMENTS", config.dummyVar),
           environ.Variable("DEPRECATED_DC_METRICS_TAG", cluster.metricsDC),
           environ.Variable("DEPRECATED_SERVER_TYPE_METRICS_TAG", "job-manager"),
       ]) +
       container.VolumeMounts([
           // TODO: Move the sidecars to mixins!
           logs.analyticsVolumeMount("/var/analytics"),
           logs.metricsVolumeMount("/var/metrics"),
           logs.serviceVolumeMount("/var/log/service"),
           pki.volumeMount(),
           // appconfd.volumeMount(),
           // logback.volumeMount(),
       ]) +
       container.Resources(config.resourceLimits) + {
           securityContext: securityContext.defaultCapabilities(),
       };
}
```

And your teammate could write the following code to define 
the VolumeMounts:

```javascript
local core = import "../../kube/core.libsonnet";
local util = import "../../kube/util.libsonnet";

Sidecar(containerNames)::
       local containerNameSet = std.set(containerNames);
       deployment.MapContainers(
           function(podContainer)
               if std.length(std.setInter([podContainer.name], containerNameSet)) > 0
               then podContainer + container.VolumeMounts([self.volumeMount()])
               else podContainer
       ) +
       deployment.mixin.podTemplate.Volumes([self.volume()]),
```

Then, you write a deployment definition for the container that 
adds the VolumeMounts for logging using mixins:

```javascript
local core = import "../../kube/core.libsonnet";
local util = import "../../kube/util.libsonnet";

local deployment = util.app.v1beta1.deployment;

{
  deployment.FromContainer( // deployment
       config.deploymentName,
       config.replicas,
       appContainer,
       podLabels=appLabels) +
   deployment.mixin.metadata.Namespace(namespace.name) +
   deployment.mixin.spec.RevisionHistoryLimit(2) +
   deployment.mixin.spec.Selector(appLabels) +
   deployment.mixin.podTemplate.NodeSelector({
       "box.com/pool": config.poolLabel,
   }) +
   deployment.mixin.podTemplate.Containers(logtailer.containers) +
   deployment.mixin.podTemplate.Volumes(logtailer.volumes + [
       logs.analyticsVolume(namespace.name),
       logs.serviceVolume(namespace.name),
       pki.volume(),
       logtailer.metricsVolume(namespace.name),
       // appconfd.volume(),
       // logback.volume(),
   ]) +
   logback.Sidecar([app.name]); // sidecar added
}
```

[readme]: ../readme.md "ksonnet readme"


