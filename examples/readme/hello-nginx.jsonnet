// This expects to be run with `jsonnet -J <path to ksonnet-lib>`
local k = import "ksonnet.beta.1/k.libsonnet";
local util = import "ksonnet.beta.1/util.libsonnet";

local container = k.core.v1.container;
local deployment = k.apps.v1beta1.deployment;

local nginxContainer =
  container.default("nginx", "nginx:1.7.9") +
  container.helpers.namedPort("http", 80);

util.prune(
  deployment.default("nginx-deployment", nginxContainer) +
    deployment.mixin.spec.template({metadata: {labels: {app: "nginx"}}}) +
    deployment.mixin.spec.replicas(2))
