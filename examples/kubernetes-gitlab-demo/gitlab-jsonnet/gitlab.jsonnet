// Configures GitLab to run on Kubernetes, including storage, deployments, and
// services.

local core = import "../../../kube/core.libsonnet";
local kubeUtil = import "../../../kube/util.libsonnet";

local gitlab = import "./gitlab.libsonnet";
local postgres = import "./postgres.libsonnet";
local redis = import "./redis.libsonnet";
local config = import "./config.libsonnet";

//
// Configuration.
//

local gitlabDeploymentName = "gitlab";
local postgresDeploymentName = "gitlab-postgresql";
local redisDeploymentName = "gitlab-redis";

// GitLab app configurator.
local configurator = config.Configurator("gitlab");

//
// Application.
//

{
  "gitlab-ns.json": core.v1.namespace.Default(configurator.namespace),

  // Gitlab app services and deployments.
  "gitlab-deployment.json": gitlab.Deployment(
    configurator, gitlabDeploymentName, gitlabDeploymentName),
  "postgres-deployment.json": postgres.Deployment(
    configurator, postgresDeploymentName, postgresDeploymentName),
  "redis-deployment.json": redis.Deployment(
    configurator, redisDeploymentName, redisDeploymentName),

  "gitlab-svc.json": gitlab.Service(
    configurator, gitlabDeploymentName, gitlabDeploymentName),
  "postgres-svc.json": postgres.Service(
    configurator, postgresDeploymentName, postgresDeploymentName),
  "redis-svc.json": redis.Service(
    configurator, redisDeploymentName, redisDeploymentName),

  // Gitlab app config maps.
  "gitlab-config.json": gitlab.AppConfigMap(configurator),
  "gitlab-build-patches.json": gitlab.PatchesConfigMap(configurator),
  "postgresql-configmap.json": postgres.ConfigMap(configurator),

  // GitLab app secrets.
  "gitlab-secrets.json": gitlab.AppSecrets(configurator),

  // GitLab app volumes.
  "gitlab-config-storage.json": gitlab.ConfigStorageClaim(configurator),
  "gitlab-rails-storage.json": gitlab.RailsStorageClaim(configurator),
  "gitlab-registry-storage.json": gitlab.RegistryStorageClaim(configurator),
  "postgresql-storage.json:": postgres.StorageClaim(configurator),
  "redis-storage.json": redis.StorageClaim(configurator),
}