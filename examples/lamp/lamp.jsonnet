local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local container = core.v1.container;
local claim = core.v1.volume.claim;
local pod = core.v1.pod;
local deployment = kubeUtil.app.v1beta1.deployment;
local persistent = core.v1.volume.persistent;
local mount = core.v1.volume.mount;

{
  local lampPvcName = "my-lamp-site-data",
  local configStorageVolume =
    persistent.Default("site-data", lampPvcName),

  local mysqlContainer =
    container.Default("mysql", "mysql") +
    container.VolumeMounts([
      mount.FromVolume(configStorageVolume, "/var/lib/mysql") {
        subPath: "mysql"
      },
    ]),


  local phpContainer =
    container.Default("php", "php") +
    container.NamedPort("http", 80) +
    container.NamedPort("https", 443) +
    container.VolumeMounts([
      mount.FromVolume(configStorageVolume, "/var/www/html") {
        subPath: "html"
      },
    ]),

  "lamp-pod.json": kubeUtil.app.v1.pod.FromContainer(mysqlContainer) {
      spec+: {
        containers+: [phpContainer],
      } +
      pod.spec.Volumes([configStorageVolume]),
    },

  // NOTE: This is a total guess.
  "lamp-pvc.json": claim.DefaultPersistent(
    "lamp-test", lampPvcName, ["ReadWriteMany"], "20Gi"),
}
