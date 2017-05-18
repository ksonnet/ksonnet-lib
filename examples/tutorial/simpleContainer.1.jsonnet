local kubeCore = import "../../kube/core.libsonnet";
local container = kubeCore.v1.container;

container.Default("nginx", "nginx:1.7.9") +
container.NamedPort("http", 80)