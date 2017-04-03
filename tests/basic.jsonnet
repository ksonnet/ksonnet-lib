local v1 = (import "../kube/core.libsonnet").v1;
local v1beta1 = (import "../kube/core.libsonnet").extensions.v1beta1;

{
  claimTest1:
    v1.volume.claim.DefaultPersistent(
      "pvcName1", ["ReadWrite"], "2G", namespace="pvcNamespace1"),
  configMapTest1:
    v1.configMap.Default("namespace1", "configMap1", {datum1: "value1"}),
  containerPortTest1:
    v1.port.container.Default(9090),
  containerPortTest2:
    v1.port.container.Named("containerPort1", 9091),
  containerPortTest3:
    v1.port.container.Default(9092) +
    v1.port.container.Name("containerPort2") +
    v1.port.container.Protocol("UDP") +
    v1.port.container.HostPort(9093) +
    v1.port.container.HostIp("127.0.0.2"),
  deploymentTest1:
    local deployment = v1beta1.deployment;
    deployment.Default("hello", {}) +
    deployment.mixin.spec.Selector({frog: "ribbit"}) +
    deployment.mixin.spec.MinReadySeconds(3) +
    deployment.mixin.spec.RollingUpdateStrategy(),
  metadataTest1:
    v1.metadata.Default() +
    v1.metadata.Annotations({ annotation1: "label1" }) +
    v1.metadata.Annotations({ annotation2: "label2" }) +
    v1.metadata.Annotation("annotation3", "label3") +
    v1.metadata.Name("name1") +
    v1.metadata.Labels({label1: "value1"}) +
    v1.metadata.Labels({label2: "value2"}) +
    v1.metadata.Label("label3", "value3") +
    v1.metadata.Namespace("namespace1"),
  mountTest1:
    v1.volume.mount.Default("volume1", "/path/to/mount1"),
  mountTest2:
    v1.volume.mount.FromVolume(
      v1.volume.persistent.Default("mountPv1", "mountPvc1"),
      "/path/to/mount2"),
  mountTest3:
    v1.volume.mount.FromConfigMap(
      v1.volume.configMap.Default("mountConfigMap1", "mountConfigMapName1"),
      "/path/to/mount2"),
  namespaceTest1:
    v1.namespace.Default("namespace1"),
  probeTest1:
    v1.probe.Default(1),
  probeTest2:
    v1.probe.Http("/probePath1", "probePort1", 3),
  probeTest3:
    v1.probe.Tcp(33, 4),
  probeTest4:
    v1.probe.Exec("execCommand1", 5),
  secretTest1:
    v1.secret.Default("namespace1", "secret1", {secretKey1: "secretValue1"}) +
    v1.secret.StringData("data1") +
    v1.secret.Type("type1"),
  servicePortTest1:
    v1.port.service.Default(8080),
  servicePortTest2:
    v1.port.service.WithTarget(8081, 8082),
  servicePortTest3:
    v1.port.service.Named("servicePort1", 8083, 8084),
  servicePortTest4:
    v1.port.service.Default(8085) +
    v1.port.service.Name("servicePort2") +
    v1.port.service.Protocol("TCP") +
    v1.port.service.TargetPort(8086) +
    v1.port.service.NodePort(8087),
  serviceTest1:
    v1.service.Default("service1", []) +
    v1.service.Metadata(v1.metadata.Name("serviceName1")) +
    v1.service.mixin.metadata.Namespace("namespace1") +
    v1.service.Spec(v1.service.spec.SessionAffinity("ClientIP")) +
    v1.service.mixin.spec.ExternalName("externalName1"),
  volumeConfigMap1:
    v1.volume.configMap.Default("configMap1", "configMapName1"),
  volumeEmptyDir1:
    v1.volume.emptyDir.Default("emptyDir1"),
  volumeHostPathTest1:
    v1.volume.hostPath.Default("hostPath1", "/path/to/nowhere"),
  volumePersistentTest1:
    v1.volume.persistent.Default("pv1", "pvc1"),
  volumePersistentTest2:
    v1.volume.persistent.DefaultFromClaim(
      "pv2",
      v1.volume.claim.DefaultPersistent(
        "pvc2", ["ReadWrite"], "2G", namespace="pvNamespace1")),
  volumeSecretTest1:
    v1.volume.secret.Default("secretVolume1", "secretVolumeName1"),
}