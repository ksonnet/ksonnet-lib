//
// NOTE: You will need Jsonnet PR #318 in order to compile this code.
//

local k = import "k.libsonnet";
local deployment = k.extensions.v1beta1.deployment;
local container = deployment.mixin.spec.template.spec.containersType;
local containerPort = container.portsType;
local service = k.core.v1.service;
local servicePort = service.mixin.spec.portsType;

// Common parts.
local commonLabels = { app: "hello", tier: "backend" };
local conatinerPortName = "http";

// Create nginx container with container port 80 open.
local helloGke =
  container.new("nginx", "nginx:1.13.0") +
  container.ports(containerPort.named(conatinerPortName, 80));

k.core.v1.list.new([
  // Create default `Deployment` object from nginx container.
  deployment.new(
    "nginx-deployment", 3, helloGke, commonLabels + {track: "stable"}),

  // Create default `Deployment` object from nginx container.
  service.new("hello", commonLabels, [
    servicePort.tcp(80, conatinerPortName)
  ]),
])
