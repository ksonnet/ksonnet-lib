{
  local ApiVersion = { apiVersion: "v1", },

  Default(name, version, description, maintainers, engine="gotpl")::
    ApiVersion +
    self.Name(name) +
    self.Version(version) +
    self.Description(description) +
    self.Maintainers(maintainers) +
    self.Engine(engine),

  Description(text):: { description: text, },
  Keywords(keywordsList):: { keywords: keywordsList, },
  Home(homeUrl):: { home: homeUrl, },
  Icon(url):: { icon: url, },
  Name(name):: { name: name, },
  Version(version):: { version: version, },
  Sources(urlList):: { sources: urlList, },
  Maintainers(maintainerList):: { maintainers: maintainerList, },
  Engine(engine):: { engine: engine, },

  metadata:: {
    Default(chart, fullname, name, releaseName, releaseService):: {
      labels: {
        app: std.extVar("name"),
        chart: "%s-%s" % [chart.name, chart.version],
        heritage: {{ .Release.Service }},
        release: {{ .Release.Name }},
      },
      name: std.extVar("fullname"),
    },
  },

  maintainer:: {
    Default(name, email):: { name: name, email: email, },
  }
}