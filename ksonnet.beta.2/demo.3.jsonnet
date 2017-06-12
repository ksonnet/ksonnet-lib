//
// NOTE: You will need Jsonnet PR #318 in order to compile this code.
//

local k = import "k.libsonnet";
local deployment = k.extensions.v1beta1.deployment;
local container = deployment.mixin.spec.template.spec.containersType;
local containerPort = container.portsType;
local mount = container.volumeMountsType;
local volume = deployment.mixin.spec.template.spec.volumesType;

// Another time writes this.
local appendVolumes(path) =
  local nginxMount = mount.new("mypd", path);
  local nginxVol = volume.fromPvc("mypd", "myclaim-1");
  deployment.mapContainers(
    function(c) c + container.volumeMounts(nginxMount)) +
  deployment.mixin.spec.template.spec.volumes(nginxVol);

// Your team writes the app code.
local helloGke =
  container.new("nginx", "nginx:1.13.0") +
  container.ports(containerPort.named("http", 80));
deployment.new(
  "nginx-deployment", 3, helloGke, {app: "nginx"}) +
appendVolumes("/usr/share/nginx/html")
