local chart = import "../chart.libsonnet";
local maintainer = chart.maintainer;

local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";
local v1 = core.v1;
local v1beta1 = core.extensions.v1beta1;

local values = import "values.libsonnet";

{
  // Configuration.

  local chartSpec =
    chart.Default("chaoskube", "0.5.0") +
    chart.Description(
      "Chaoskube periodically kills random pods in your Kubernetes cluster.") +
    chart.Home("https://github.com/linki/chaoskube") +
    chart.Source("https://github.com/linki/chaoskube") +
    chart.Maintainer(
      maintainer.Default("Martin Linkhorst", "linki+kubernetes.io@posteo.de")),

  "chart.json": chartSpec,
  "deployment.json":
    local deployment = v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
    local container = v1.container;
    local appContainer =
      container.Default(values.name, containerImage) +
      container.Args([
        "--in-cluster",
        "--interval=%s" % values.interval,
        "--labels=%s" % values.labels,
        "--annotations=%s" % values.annotations,
        "--namespaces=%s" % values.namespaces,
      ] + if !values.dryRun then ["--no-dry-run"] else []) +
      container.Resources(
        container.resources.Requests(
          values.resources.cpu, values.resources.memory) +
        container.resources.Limits(
          values.resources.cpu, values.resources.memory));
    local labels = chart.DefaultLabels(fullname, chartSpec, release);
    deployment.FromContainer(
      fullname, values.replicas, appContainer, labels=labels, podLabels=labels),

  // Data

  local containerImage = chart.ContainerImage(values.image, values.imageTag),
  local release = {
    name: std.extVar("release.name"),
    service: std.extVar("release.service"),
  },
  local fullname = chart.Fullname(release.name, release.service),
}