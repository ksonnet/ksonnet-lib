//
// NOTE: You will need Jsonnet PR #318 in order to compile this code.
//

local k = import "k.libsonnet";
local deployment = k.extensions.v1beta1.deployment;
local container = deployment.mixin.spec.template.spec.containersType;
local containerPort = container.portsType;

// Create nginx container with container port 80 open.
local nginxContainer =
  container.new("nginx", "nginx:1.13.0") +
  container.ports(containerPort.named("http", 80));

// Create default `Deployment` object from nginx container.
deployment.new(
  "nginx-deployment", 3, nginxContainer, {app: "nginx"})
