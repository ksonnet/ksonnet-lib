local core = import "../kube/core.libsonnet";

{
  metadataTest:
    core.v1.metadata.Default() +
    core.v1.metadata.Annotations({ cow: "moo" }) +
    core.v1.metadata.Annotations({ chicken: "cluck" }),
  deployment:
    local deployment = core.extensions.v1beta1.deployment;
    deployment.Default("hello", {}) +
    deployment.mixin.spec.Selector({frog: "ribbit"}) +
    deployment.mixin.spec.MinReadySeconds(3) +
    deployment.mixin.spec.RollingUpdateStrategy(),
}