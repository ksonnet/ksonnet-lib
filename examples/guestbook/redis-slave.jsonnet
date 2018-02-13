local k = import "../../ksonnet.beta.1/k.libsonnet";
local util = import "../../ksonnet.beta.1/util.libsonnet";

local rc = k.core.v1.replicationController;
local container = k.core.v1.container;
local service = k.core.v1.service;

local redisContainer =
  container.default("redis-slave", "kubernetes/redis-slave:v2") +
  container.ports([{name: "redis-server", containerPort: 6379}]);

local redisLabel = {app: "redis", role: "slave"};

local redisSlaveRC =
  rc.default("redis-slave") +
  rc.mixin.spec.replicas(2) +
  {spec+: {selector+: redisLabel}} +
  {spec+: {template+: {spec+: {containers: [redisContainer]}}}} +
  {spec+: {template+: {metadata: {labels: redisLabel}}}};

local redisSlaveSvc =
  service.default("redis-slave") +
  service.mixin.spec.ports([{port: 6379, targetPort: "redis-server"}]) +
  service.mixin.spec.selector(redisLabel) +
  service.mixin.metadata.labels(redisLabel);

{
  redisSlaveRC : util.prune(redisSlaveRC),
  redisSlaveSvc: util.prune(redisSlaveSvc),
}
