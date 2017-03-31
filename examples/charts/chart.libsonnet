local core = import "../../kube/core.libsonnet";
local kubeUtil = import "../../kube/util.libsonnet";

local service = core.v1.service;
local deployment = core.extensions.v1beta1.deployment;

{
  local apiVersion = { apiVersion: "v1" },

  local dnsName(name) =
    if std.length(name) <= 63
    then name
    else std.substr(name, 0, 63),

  ContainerImage(image, tag):: "%s:%s" % [image, tag],

  Name(chartName)::
    // TODO: Remove trailing '-' character if it exists after truncation.
    dnsName(chartName),

  Fullname(chartName, releaseName)::
    local name = "%s-%s" % [chartName, releaseName];
    dnsName(name),

  DefaultLabels(name, chart, release):: {
    "app": name,
    "chart": chart.name,
    "heritage": release.service,
    "release": release.name,
  },

  DefaultSelector(name, release)::
    service.mixin.spec.Selector({
      app: name,
      release: release.name
    }),

  DefaultService(fullname, name, chartSpec, release)::
    service.Default(fullname, []) +
    service.mixin.metadata.Namespace(fullname) +
    service.Metadata(
      core.v1.metadata.Labels(self.DefaultLabels(name, chartSpec, release))) +
    self.DefaultSelector(name, release),

  Default(name, version, engine="gotpl"):: apiVersion + {
    name: name,
    version: version,
    engine: engine,
    sources: [],
    maintainers: [],
  },

  Description(description):: { description: description },
  Keywords(keywordsList):: { keywords: keywordsList, },
  Home(homeUrl):: { home: homeUrl, },
  Icon(url):: { icon: url, },
  Version(version):: { version: version, },
  Source(url):: { sources+: [url] },
  Sources(urlList):: { sources+: urlList, },
  Maintainer(maintainer):: { maintainers+: [maintainer] },
  Maintainers(maintainerList):: { maintainers+: maintainerList, },
  Engine(engine):: { engine: engine, },

  maintainer:: {
    Default(name, email):: { name: name, email: email },
  },
}