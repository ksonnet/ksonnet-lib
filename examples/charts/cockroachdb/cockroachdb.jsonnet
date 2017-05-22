// This file rewrites the Cockroachdb Helm chart in ksonnet.
// Compare with https://github.com/kubernetes/charts/blob/master/stable/cockroachdb/templates/cockroachdb-petset.yaml

local chart = import "../chart.libsonnet";
local maintainer = chart.maintainer;

local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";
local v1 = core.v1;
local v1beta1 = core.extensions.v1beta1 + core.policy.v1beta1 + core.apps.v1beta1;

local values = import "values.libsonnet";

local claim = v1.volume.claim;
local container = v1.container;
local distruptionBudget = v1beta1.podDistruptionBudget;
local env = v1.env;
local statefulSet = v1beta1.statefulSet;
local mount = v1.volume.mount;
local pod = v1.pod;
local port = v1.port;
local probe = v1.probe;
local service = v1.service;
local volume = v1.volume;

{
  //
  // Chart specification.
  //

  "chart.json": chartSpec,

  local chartSpec =
    chart.Default("cockroachdb", "0.2.2") +
    chart.Description("CockroachDB is a scalable, survivable, strongly-consistent SQL database.") +
    chart.Home("https://www.cockroachlabs.com") +
    chart.Icon("https://raw.githubusercontent.com/cockroachdb/cockroach/master/docs/media/cockroach_db.png") +
    chart.Source("https://github.com/cockroachdb/cockroach") +
    chart.Maintainer(
      maintainer.Default("Alex Robinson", "alex@cockroachlabs.com")),

  //
  // Services.
  //

  // This service is meant to be used by clients of the database. It
  // exposes a ClusterIP that will automatically load balance
  // connections to the different database pods.
  "publicService.json": publicService,

  // This service only exists to create DNS entries for each pod in
  // the stateful set such that they can resolve each other's IP
  // addresses. It does not create a load-balanced ClusterIP and
  // should not be used directly by clients in most circumstances.
  "privateService.json":
    publicService +
    service.mixin.metadata.Name(privateName) +
    // This is needed to make the peer-finder work properly and to
    // help avoid edge cases where instance 0 comes up after losing
    // its data and needs to decide whether it should create a new
    // cluster or try to join an existing one. If it creates a new
    // cluster when it should have joined an existing one, we'd end
    // up with two separate clusters listening at the same service
    // endpoint, which would be very bad.
    service.mixin.metadata.annotation.TolerateUnreadyEndpoints(true) +
    addon.service.prometheus +
    service.mixin.spec.ClusterIp("None"),

  local publicService =
    service.Default("%s-public" % privateName, servicePorts, labels=labels) +
    service.mixin.spec.Selector(componentLabel),

  //
  // Policies.
  //

  "disruptionPolicy.json":
    distruptionBudget.Default("%s-budget" % privateName, labels) +
    distruptionBudget.mixin.spec.Selector(componentLabel) +
    distruptionBudget.mixin.spec.MinAvailable(values.MinAvailable),

  //
  // Stateful sets.
  //

  "statefulSet.json":
    local dataVolumeClaimName = "datadir";
    local dataVolume =
      volume.persistent.Default("datadir", dataVolumeClaimName);
    local dataVolumeMount =
      mount.FromVolume(dataVolume, "/cockroach/cockroach-data");
    local dataVolumeClaim =
      claim.DefaultPersistent(
        dataVolumeClaimName, ["ReadWriteOnce"], values.Storage) +
      claim.mixin.metadata.annotation.AlphaStorageClass(values.StorageClass);
    local httpProbe = probe.Http("/_admin/v1/health", "http", 30);
    local appContainer =
      container.Default(privateName, containerImage, values.ImagePullPolicy) +
      container.Ports(containerPorts) +
      container.Resources(values.Resources) +
      container.Env([env.Variable("STATEFULSET_NAME", privateName)]) +
      container.LivenessProbe(httpProbe) +
      container.ReadinessProbe(httpProbe) +
      container.VolumeMounts([dataVolumeMount]) +
      container.Command(containerCommand);
    local podTemplate =
      kubeUtil.app.v1.pod.template.FromContainer(
        appContainer, labels=labels, volumes=[dataVolume]) +
      pod.template.mixin.metadata.annotation.PodAffinity(podAffinitySpec) +
      pod.template.mixin.metadata.annotation.PodInitContainers(podInitSpec);
    statefulSet.Default(
      privateName, values.Replicas, podTemplate, volumeClaimTemplates=[dataVolumeClaim]),

  //
  // Addons.
  //

  local addon = {
    service:: {
      prometheus:: service.mixin.metadata.Annotations({
        // Enable automatic monitoring of all instances when Prometheus
        // is running in the cluster.
        "prometheus.io/scrape": true,
        "prometheus.io/path": "_status/vars",
        "prometheus.io/port": "8080",
      }),
    },
  },

  //
  // Data.
  //

  local release = {
    name: std.extVar("release.name"),
    service: std.extVar("release.service"),
  },
  local containerImage = chart.ContainerImage(values.Image, values.ImageTag),
  local privateName = chart.Fullname(release.name, values.Name, truncLength=56),
  local componentName = "%s-%s" % [release.name, values.Component],
  local componentLabel = {component: componentName},
  // TODO: We need to remove the `app` label, which comes by default.
  local labels =
    chart.DefaultLabels(privateName, chartSpec, release) +
    componentLabel,
  local servicePorts = [
    // The main port, served by gRPC, serves Postgres-flavor SQL,
    // internode traffic and the cli.
    port.service.Named("grpc", values.GrpcPort, values.GrpcPort),
    // The secondary port serves the UI as well as health and debug
    // endpoints.
    port.service.Named("http", values.HttpPort, values.HttpPort),
  ],
  local containerPorts = [
    port.container.Named("grpc", values.GrpcPort),
    port.container.Named("http", values.HttpPort),
  ],
  local containerCommand = [
    "/bin/bash",
    "-ecx",
|||
    # The use of qualified `hostname -f` is crucial:
    # Other nodes aren't able to look up the unqualified hostname.
    CRARGS=("start" "--logtostderr" "--insecure" "--host" "$(hostname -f)" "--http-host" "0.0.0.0")
    # We only want to initialize a new cluster (by omitting the join flag)
    # if we're sure that we're the first node (i.e. index 0) and that
    # there aren't any other nodes running as part of the cluster that
    # this is supposed to be a part of (which indicates that a cluster
    # already exists and we should make sure not to create a new one).
    # It's fine to run without --join on a restart if there aren't any
    # other nodes.
    if [ ! "$(hostname)" == "${STATEFULSET_NAME}-0" ] || \
        [ -e "/cockroach/cockroach-data/cluster_exists_marker" ]
    then
      CRARGS+=("--join" "${STATEFULSET_NAME}-public")
    fi
    exec /cockroach/cockroach ${CRARGS[*]}
|||
  ],
  local podAffinitySpec =
|||
    {
      "podAntiAffinity": {
        "preferredDuringSchedulingIgnoredDuringExecution": [{
          "weight": 100,
          "labelSelector": {
            "matchExpressions": [{
              "key": "component",
              "operator": "In",
              "values": ["{{.Release.Name}}-{{.Values.Component}}"]
            }]
          },
          "topologyKey": "kubernetes.io/hostname"
        }]
      }
    }
|||,
  // Init containers are run only once in the lifetime of a pod, before
  // it's started up for the first time. It has to exit successfully
  // before the pod's main containers are allowed to start.
  // This particular init container does a DNS lookup for other pods in
  // the set to help determine whether or not a cluster already exists.
  // If any other pods exist, it creates a file in the cockroach-data
  // directory to pass that information along to the primary container that
  // has to decide what command-line flags to use when starting CockroachDB.
  // This only matters when a pod's persistent volume is empty - if it has
  // data from a previous execution, that data will always be used.
  // The cockroachdb/cockroach-k8s-init image is defined at
  // github.com/cockroachdb/cockroach/blob/master/cloud/kubernetes/init
  local podInitSpec =
|||
    [
        {
            "name": "bootstrap",
            "image": "{{.Values.BootstrapImage}}:{{.Values.BootstrapImageTag}}",
            "imagePullPolicy": "{{.Values.ImagePullPolicy}}",
            "args": [
              "-on-start=/on-start.sh",
              "-service={{ printf "%s-%s" .Release.Name .Values.Name | trunc 56 }}"
            ],
            "env": [
              {
                  "name": "POD_NAMESPACE",
                  "valueFrom": {
                      "fieldRef": {
                          "apiVersion": "v1",
                          "fieldPath": "metadata.namespace"
                      }
                  }
                }
            ],
            "volumeMounts": [
                {
                    "name": "datadir",
                    "mountPath": "/cockroach/cockroach-data"
                }
            ]
        }
    ]
|||,
}
