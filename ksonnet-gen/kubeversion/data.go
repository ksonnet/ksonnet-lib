package kubeversion

//-----------------------------------------------------------------------------
// Kubernetes version-specific data for customizing code that's
// emitted.
//-----------------------------------------------------------------------------

var versions = versionSet{
	"v1.7.0": versionData{
		idAliases: idAliasSet{
			"hostIPC":                        "hostIpc",
			"hostPID":                        "hostPid",
			"targetCPUUtilizationPercentage": "targetCpuUtilizationPercentage",
			"externalID":                     "externalId",
			"podCIDR":                        "podCidr",
			"providerID":                     "providerId",
			"bootID":                         "bootId",
			"machineID":                      "machineId",
			"systemUUID":                     "systemUuid",
			"volumeID":                       "volumeId",
			"diskURI":                        "diskUri",
			"targetWWNs":                     "targetWwns",
			"datasetUUID":                    "datasetUuid",
			"pdID":                           "pdId",
			"scaleIO":                        "scaleIo",
			"podIP":                          "podIp",
			"hostIP":                         "hostIp",
			"clusterIP":                      "clusterIp",
			"externalIPs":                    "externalIps",
			"loadBalancerIP":                 "loadBalancerIp",
		},
		propertyBlacklist: blackList{
			"io.k8s.kubernetes.pkg.apis.apps.v1beta1.Deployment": newPropertySet("status"),
		},
	},
}
