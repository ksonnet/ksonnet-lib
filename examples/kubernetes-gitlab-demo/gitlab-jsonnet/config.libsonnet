{
  //
  // GitLab app configurator.
  //

  Configurator(
    namespace,
    appConfigMapName="gitlab-config",
    appSecretName="gitlab-secrets",
    appPatchesConfigMapName="gitlab-patches",
    postgresConfigMapName="gitlab-postgresql-initdb",
    appConfigStorageClaimName="gitlab-config-storage",
    appDataClaimName="gitlab-rails-storage",
    appRegistryClaimName="gitlab-registry-storage",
    postgresStorageClaimName="gitlab-postgresql-storage",
    redisStorageClaimName="gitlab-redis-storage",
    registryServicePort=8105,
    mattermostServicePort=8065,
    workhorseServicePort=8005,
    sshServicePort=22,
    prometheusServicePort=9090,
    nodeExporterServicePort=9100,
    postgresServicePort=5432,
    redisServicePort=6379,
  ):: {
    namespace: namespace,

    // Config map names.
    appConfigMapName: appConfigMapName,
    appSecretName: appSecretName,
    appPatchesConfigMapName: appPatchesConfigMapName,
    postgresConfigMapName: postgresConfigMapName,

    // PVC names.
    appConfigStorageClaimName: appConfigStorageClaimName,
    appDataClaimName: appDataClaimName,
    appRegistryClaimName: appRegistryClaimName,
    postgresStorageClaimName: postgresStorageClaimName,
    redisStorageClaimName: redisStorageClaimName,

    // Service ports.
    registryServicePort: registryServicePort,
    mattermostServicePort: mattermostServicePort,
    workhorseServicePort: workhorseServicePort,
    sshServicePort: sshServicePort,
    prometheusServicePort: prometheusServicePort,
    "node-exporterServicePort": nodeExporterServicePort,
    postgresServicePort: postgresServicePort,
    redisServicePort: redisServicePort,
  },
}