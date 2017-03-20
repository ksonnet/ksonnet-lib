{
  autoscalingGroups: [],
  extraArgs: {},
  resources: {},
  podAnnotations: {},
  nodeSelector: {},
  replicaCount: 1,
  image: {
    repository: "gcr.io/google_containers/cluster-autoscaler",
    tag: "v0.4.0",
    pullPolicy: "IfNotPresent",
  },
  awsRegion: "us-east-1",
  scaleDownDelay: "10m",
  skipNodes: {
    withLocalStorage: false,
    withSystemPods: true,
  },
  service: {
    annotations: {},
    clusterIp: "",

    ## List of IP addresses at which the service is available
    ## Ref: https://kubernetes.io/docs/user-guide/services/#external-ips
    ##
    externalIps: [],

    loadBalancerIp: "",
    loadBalancerSourceRanges: [],
    servicePort: 8085,
    type: "ClusterIP",
  },
}