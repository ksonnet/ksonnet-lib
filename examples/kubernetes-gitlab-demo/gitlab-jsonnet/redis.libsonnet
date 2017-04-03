local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";

// Convenient namespaces.
local claim = core.v1.volume.claim;
local container = core.v1.container;
local deployment = core.extensions.v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
local metadata = core.v1.metadata;
local persistent = core.v1.volume.persistent;
local pod = core.v1.pod;
local port = core.v1.port + kubeUtil.app.v1.port;
local probe = core.v1.probe;
local mount = core.v1.volume.mount;
local service = core.v1.service;

{
  //
  // Deployment.
  //

  Deployment(config, deploymentName, podName)::
    // Volumes.
    local dataVolume =
      persistent.Default("data", config.redisStorageClaimName);

    // Container and pod definitions.
    local probeCommand = ["redis-cli", "ping"];
    local redisContainer =
      container.Default("redis", "redis:3.2.4", "IfNotPresent") +
      container.LivenessProbe(probe.Exec(probeCommand, 30, 5)) +
      container.ReadinessProbe(probe.Exec(probeCommand, 5, 1)) +
      container.VolumeMounts([
        mount.FromVolume(dataVolume, "/var/lib/redis"),
      ]) +
      container.Ports(podPorts);

    // Deployment.
    deployment.FromContainer(
      deploymentName, 1, redisContainer, podLabels={ name: podName }) +
    deployment.mixin.metadata.Namespace(config.namespace) +
    deployment.mixin.podTemplate.Volumes([dataVolume]),

  //
  // Service.
  //

  Service(config, serviceName, targetPod)::
    local servicePorts = port.service.array.FromContainerPorts(
      function (containerPort) config[containerPort.name + "ServicePort"],
      podPorts);

    service.Default(serviceName, servicePorts) +
    service.Metadata(
      metadata.Namespace(config.namespace) +
      metadata.Label("name", serviceName)) +
    service.mixin.spec.Selector({ name: targetPod }),

  //
  // Persistent volume claims.
  //

  StorageClaim(config)::
    local claimName = config.redisStorageClaimName;
    claim.DefaultPersistent(
      claimName, ["ReadWriteOnce"], "5Gi", namespace=config.namespace) +
    claim.mixin.metadata.annotation.BetaStorageClass("fast"),

  //
  // Private helpers.
  //

  local podPorts = [
    port.container.Named("redis", 6379),
  ],
}