local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local container = core.v1.container;
local deployment = kubeUtil.app.v1beta1.deployment;
local mount = core.v1.volume.mount;
local persistent = core.v1.volume.persistent;

local Test() =
  local nginxContainer =
    container.Default("nginx", "nginx:1.7.9") +
    container.NamedPort("http", 80);

  deployment.FromContainer("nginx-deployment", 2, nginxContainer);

local ConfigureStorage() =
  local configStorageVolume =
    persistent.Default("config", "volumeClaim");

  container.VolumeMounts([
    mount.FromVolume(configStorageVolume, "/etc/nginx"),
  ]);

{
  "deployment.json": Test(),
}
