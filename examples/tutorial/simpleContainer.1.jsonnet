local core = import "../../kube/core.libsonnet";
local container = core.v1.container;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80)