local chart = import "../chart.libsonnet";
local maintainer = chart.maintainer;

local core = import "../../../kube/core.libsonnet";
local container = core.v1.container;
local persistent = core.v1.volume.persistent;
local mount = core.v1.volume.mount;
local port = core.v1.port;
local volume = core.v1.volume;

local kubeUtil = import "../../../kube/util.libsonnet";
local deployment = core.extensions.v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
local service = core.v1.service + kubeUtil.app.v1.service;
local podTemplate = core.v1.pod.template + kubeUtil.app.v1.pod.template;

local values = import "./values.libsonnet";

local template = import "template.libsonnet";

{
  // Configuration.

  local chartSpec =
    chart.Default("aws-cluster-autoscaler", "0.2.1") +
    chart.Description("Scales worker nodes within autoscaling groups.") +
    chart.Source("https://github.com/kubernetes/contrib/tree/master/cluster-autoscaler/cloudprovider/aws") +
    chart.Maintainer(
      maintainer.Default("Michael Goodness", "mgoodness@gmail.com")),

  "chart.json": chartSpec,
  "service.json":
    chart.DefaultService(fullname, name, chartSpec, release) +
    service.Port(port.service.WithTarget(values.service.servicePort, 8085)) +
    service.Type(values.service.type) +
    service.mixin.metadata.Annotations(values.service.annotations) +
    service.ClusterIp(values.service.clusterIp) +
    service.ExternalIps(values.service.externalIps) +
    service.LoadBalancerIp(values.service.loadBalancerIp) +
    service.LoadBalancerSourceRanges(values.service.loadBalancerSourceRanges),
  "deployment.json":
    local certPath = "/etc/ssl/certs/ca-certificates.crt";
    local certsVolume = volume.hostPath.Default("ssl-certs", certPath);
    local appContainer =
      container.Default(name, containerImage, values.image.pullPolicy) +
      container.Command(containerCommand) +
      container.Env([{"AWS_REGION": values.awsRegion}]) +
      container.Port(port.container.Default(8085)) +
      container.VolumeMounts([mount.FromVolume(certsVolume, certPath, true)]) +
      container.Resources(values.resources);
    local appPod =
      podTemplate.FromContainer(appContainer) +
      podTemplate.mixin.metadata.Label("release", release.name) +
      podTemplate.mixin.metadata.Annotations(values.podAnnotations);
    deployment.FromPodTemplate(fullname, values.replicaCount, appPod) +
    deployment.NodeSelector(values.nodeSelector),

  // Data

  local name = chart.Name(chartSpec.name),
  local fullname = chart.Fullname(chartSpec.name, "beta"),
  local release = {
    name: std.extVar("release.name"),
    service: std.extVar("release.service"),
  },

  local containerImage = "%s:%s" % [values.image.repository, values.image.tag],

  local autoscaleFlags = [
    "--nodes=%s:%s:%s" % [group.minSize, group.maxSize, group.name]
    for group in values.autoscalingGroups
  ],
  local extraArgsFlags = [
    "--%s=%s" % [key, values.extraArgs[key]],
    for key in std.objectFields(values.extraArgs)
  ],
  local containerCommand =
    [
      "./cluster-autoscaler",
      "--cloud-provider=aws"
    ] +
    autoscaleFlags + [
      "--scale-down-delay=%s" % values.scaleDownDelay,
      "--skip-nodes-with-local-storage=%s" % values.skipNodes.withLocalStorage,
      "--skip-nodes-with-system-pods=%s" % values.skipNodes.withSystemPods,
      "--v=4"
    ] +
    extraArgsFlags,
}
