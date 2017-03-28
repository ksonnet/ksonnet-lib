// Helpers for configuring core components of GitLab, including the deployment,
// service, and persistent volumes.

local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";

// Convenient namespaces.
local deployment = core.extensions.v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
local container = core.v1.container;
local claim = core.v1.volume.claim;
local configMap = core.v1.configMap;
local probe = core.v1.probe;
local pod = core.v1.pod;
local port = core.v1.port + kubeUtil.app.v1.port;
local service = core.v1.service;
local persistent = core.v1.volume.persistent;
local volume = core.v1.volume;
local mount = core.v1.volume.mount;

local data = import "./data.libsonnet";

{
  //
  // Deployment.
  //

  Deployment(config, deploymentName, podName)::
    // Volumes and config maps.
    local postgresConfigMap = volume.configMap.Default(
      "initdb", config.postgresConfigMapName);
    local dataVolume = persistent.Default(
      "data", config.postgresStorageClaimName);

    // Container and pod definition for Postgres instance.
    local probeCommand = ["pg_isready", "-h", "localhost", "-U", "postgres"];
    local postgresContainer =
      container.Default(deploymentName, "postgres:9.5.3", "IfNotPresent") +
      container.Env(data.postgres.deploy.Env(
        config.appConfigMapName, config.appSecretName)) +
      container.LivenessProbe(probe.Exec(probeCommand, 30, 5)) +
      container.ReadinessProbe(probe.Exec(probeCommand, 5, 1)) +
      container.VolumeMounts([
        mount.FromVolume(dataVolume, "/var/lib/postgresql"),
        mount.FromConfigMap(postgresConfigMap, "/docker-entrypoint-initdb.d", true)]) +
      container.Ports(podPorts);

    // Deployment.
    deployment.FromContainer(
      deploymentName, 1, postgresContainer, podLabels={ name: podName }) +
    deployment.Namespace(config.namespace) +
    deployment.mixin.podTemplate.Volumes([dataVolume, postgresConfigMap]),

  //
  // Service.
  //

  Service(config, serviceName, targetPod)::
    local servicePorts = port.service.array.FromContainerPorts(
      function (containerPort) config[containerPort.name + "ServicePort"],
      podPorts);

    service.Default(serviceName, servicePorts) +
    service.Namespace(config.namespace) +
    service.Label("name", serviceName) +
    service.Selector({ name: targetPod }),

  //
  // Config maps.
  //

  ConfigMap(config):: configMap.Default(
      config.namespace, config.postgresConfigMapName, data.postgres.config),

  //
  // Persistent volume claims.
  //

  StorageClaim(config):: claim.DefaultPersistent(
      config.namespace, config.postgresStorageClaimName, ["ReadWriteOnce"], "30Gi"),

  //
  // Private helpers.
  //

  local podPorts = [
    port.container.Named("postgres", 5432),
  ],
}