local v1 = (import "../kube/core.libsonnet").v1;
local v1beta1 = (import "../kube/core.libsonnet").extensions.v1beta1;

{
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
  deploymentTest:
    local deployment = v1beta1.deployment;
    deployment.Default("hello", {}) +
    deployment.mixin.spec.Selector({frog: "ribbit"}) +
    deployment.mixin.spec.MinReadySeconds(3) +
    deployment.mixin.spec.RollingUpdateStrategy(),
  metadataTest:
    v1.metadata.Default() +
    v1.metadata.Annotations({ annotation1: "label1" }) +
    v1.metadata.Annotations({ annotation2: "label2" }) +
    v1.metadata.Annotation("annotation3", "label3") +
    v1.metadata.Name("name1") +
    v1.metadata.Labels({label1: "value1"}) +
    v1.metadata.Labels({label2: "value2"}) +
    v1.metadata.Label("label3", "value3") +
    v1.metadata.Namespace("namespace1"),
  namespaceTest:
    v1.namespace.Default("namespace1"),
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
    v1.service.Default("service1", []),
}