local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";

// Convenient namespaces.
local claim = core.v1.volume.claim;
local container = core.v1.container;
local deployment = core.extensions.v1beta1.deployment;
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

    local redisPodTemplate =
      pod.template.Default(
        pod.spec.Containers([redisContainer]) +
          pod.spec.Volumes([dataVolume])) +
      pod.template.Labels({ name: podName });

    // Deployment.
    deployment.Default(
      deploymentName,
      deployment.spec.ReplicatedPod(1, redisPodTemplate)) +
    deployment.Namespace(config.namespace),

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
  // Persistent volume claims.
  //

  StorageClaim(config):: claim.DefaultPersistent(
    config.namespace, config.redisStorageClaimName, ["ReadWriteOnce"], "5Gi"),

  //
  // Private helpers.
  //

  local podPorts = [
    port.container.Named("redis", 6379),
  ],
}