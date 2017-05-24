local k = import "../../ksonnet.beta.1/k.libsonnet";

local container = k.core.v1.container;
local deployment = k.apps.v1beta1.deployment;

local nginxContainer =
  container.default("nginx", "nginx:1.7.9") +
  container.helpers.namedPort("http", 80);

deployment.default("nginx-deployment", nginxContainer) +
deployment.mixin.spec.replicas(2)
