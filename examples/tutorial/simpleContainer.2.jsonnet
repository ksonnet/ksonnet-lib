local core = import "../../kube/core.libsonnet";
local container = core.v1.container;
local probe = core.v1.probe;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80) +
container.LivenessProbe(probe.Http("/", 80, 15, 1))