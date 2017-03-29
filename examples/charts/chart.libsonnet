local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local service = core.v1.service + kubeUtil.app.v1.service;
local deployment = core.extensions.v1beta1.deployment;

{
  local apiVersion = { apiVersion: "v1" },

  local dnsName(name) =
    if std.length(name) <= 63
    then name
    else std.substr(name, 0, 63),

  Name(chartName)::
    // TODO: Remove trailing '-' character if it exists after truncation.
    dnsName(chartName),

  Fullname(chartName, releaseName)::
    local name = "%s-%s" % [chartName, releaseName];
    dnsName(name),

  DefaultLabels(name, chart, release)::
    service.mixin.metadata.Label("app", name) +
    service.mixin.metadata.Label("chart", chart.name) +
    service.mixin.metadata.Label("heritage", release.service) +
    service.mixin.metadata.Label("release", release.name),

  DefaultSelector(name, release)::
    service.Selector({
      app: name,
      release: release.name
    }),

  DefaultService(fullname, name, chartSpec, release)::
    service.Default(fullname, []) +
    service.mixin.metadata.Namespace(fullname) +
    self.DefaultLabels(name, chartSpec, release) +
    self.DefaultSelector(name, release),

  Default(name, version, engine="gotpl"):: apiVersion + {
    name: name,
    version: version,
    engine: engine,
    sources: [],
    maintainers: [],
  },

  Description(description):: { description: description },

  Source(url):: { sources+: [url] },

  Maintainer(maintainer):: { maintainers+: [maintainer] },

  maintainer:: {
    Default(name, email):: { name: name, email: email },
  },
}