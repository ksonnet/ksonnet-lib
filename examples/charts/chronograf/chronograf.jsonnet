local chart = import "../chart.libsonnet";
local maintainer = chart.maintainer;

local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";
local v1 = core.v1;
local v1beta1 = core.extensions.v1beta1;

local values = import "values.libsonnet";

local deployment = v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
local container = v1.container;
local ingress = v1beta1.ingress;
local port = v1.port;
local probe = v1.probe;
local service = v1.service;
local volume = v1.volume;

{
  // Configuration.

  local chartSpec =
    chart.Default("chronograf", "0.1.2") +
    chart.Description("Open-source web application written in Go and React.js that provides the tools to visualize your monitoring data and easily create alerting and automation rules.") +
    chart.Home("https://www.influxdata.com/time-series-platform/chronograf/") +
    chart.Source("https://github.com/linki/chaoskube") +
    chart.Maintainer(
      maintainer.Default("Jack Zampolin", "jack@influxdb.com")) +
    chart.Keywords([
      "chronograf",
      "visualizaion",
      "timeseries",
    ]),

  "chart.json": chartSpec,

  "deployment.json":
    local dataVolume =
      local volumeName = "data";
      if values.persistence.enabled
      then volume.persistent.Default(volumeName, fullname)
      else volume.emptyDir.Default(volumeName);
    local appContainer =
      container.Default(fullname, containerImage, values.image.pullPolicy) +
      container.Port(port.container.Named("api", 8888)) +
      container.LivenessProbe(probe.Http("/ping", "api", 0)) +
      container.ReadinessProbe(probe.Http("/ping", "api", 0)) +
      container.VolumeMounts([
        volume.mount.FromVolume(dataVolume, "/var/lib/chronograf")
      ]) {
        resources: values.resources,
      };

    deployment.FromContainer(fullname, 1, appContainer, labels=labels),

  "service.json":
    local ports = [port.service.WithTarget(80, 8888)];
    service.Default(fullname, ports) +
    service.mixin.metadata.Labels(labels) +
    service.mixin.spec.Selector({app: fullname}) +
    service.mixin.spec.Type(values.service.type),

  // TODO: If ingress is enabled.
  [if values.ingress.enabled then "ingress.json"]:
    ingress.Default(fullname, labels=labels) +
    ingress.mixin.metadata.Annotations(values.ingress.annotations) +
    ingress.mixin.spec.Rule(
      values.ingress.hostname,
      ingress.httpIngressPath.Default(fullname, 80)) +
    if values.ingress.tls
    then ingress.mixin.spec.Tls([values.ingress.hostname], "%s-tls" % fullname)
    else {},

  // Storage.

  [if values.persistence.enabled then "pvc.json"]:
    local mixin = volume.claim.mixin;
    local storageClassMixin =
      if std.objectHas(values.persistence, "storageClass")
      then mixin.metadata.annotation.BetaStorageClass(
        values.persistence.storageClass)
      else mixin.metadata.annotation.AlphaStorageClass("default");
    volume.claim.DefaultPersistent(
      fullname, [values.persistence.accessMode], values.persistence.size) +
    mixin.metadata.Labels(labels) +
    storageClassMixin,

  // Data

  local labels = chart.DefaultLabels(fullname, chartSpec, release),
  local containerImage = chart.ContainerImage(
    values.image.repository, values.image.tag),
  local release = {
    name: std.extVar("release.name"),
    service: std.extVar("release.service"),
  },
  local fullname = chart.Fullname(release.name, release.service),
}