local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local container = core.v1.container;
local pod = core.v1.pod;
local deployment = kubeUtil.app.v1beta1.deployment + core.extensions.v1beta1.deployment;

{
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80),

  "deployment.json":
    deployment.FromContainer("nginx-deployment", 2, nginxContainer) +
    deployment.MixinSpec(
      deployment.spec.RollingUpdateStrategy(1, 1) +
      deployment.spec.Selector({ "app": "nginx" })),
}
