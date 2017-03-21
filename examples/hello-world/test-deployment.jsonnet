local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local container = core.v1.container;
local deployment = core.extensions.v1beta1.deployment + kubeUtil.app.v1beta1.deployment;
local env = core.v1.env;
local pod = core.v1.pod + kubeUtil.app.v1.pod;
local podTemplate = core.v1.pod.template + kubeUtil.app.v1.pod.template;
local port = core.v1.port;
local probe = core.v1.probe;
local volume = core.v1.volume;

{
  local serviceName = std.extVar("serviceName"),

  local secret = volume.secret.Default("gcp-credentials", "gcp-credentials"),
  local appLabels = {
    app: serviceName,
    tier: "backend",
  },

  local appContainer =
    container.Default(
      "%s-service" % serviceName,
      "gcr.io/fkorotkov/%s-service:latest" % serviceName) +
    container.Env([
      env.Variable(
        "GOOGLE_APPLICATION_CREDENTIALS",
        "/etc/credentials/service-account-credentials.json")
    ]) +
    container.Port(port.container.Default(8080)) +
    container.LivenessProbe(probe.Http("/healthz", 8080, 1, 1, 60)) +
    container.ReadinessProbe(probe.Http("/healthz", 8080, 10, 1, 60)) +
    container.VolumeMounts([
      volume.mount.FromVolume(secret, "/etc/credentials", true)
    ]),

  "deployment.json":
    deployment.FromContainer("%s-service-deployment" % serviceName, 1,
                             appContainer, labels=appLabels, podLabels=appLabels),
}
