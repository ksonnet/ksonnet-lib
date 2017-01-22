local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local container = core.v1.container;
local probe = core.v1.probe;
local pod = kubeUtil.app.v1.pod;

{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80) +
    container.LivenessProbe(probe.Http("/", 80, 15, 1)),

  "nginxPod.json": pod.FromContainer(nginxContainer),
}