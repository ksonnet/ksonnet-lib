local k = import "../../ksonnet.beta.1/k.libsonnet";
local util = import "../../ksonnet.beta.1/util.libsonnet";

local rc = k.core.v1.replicationController;
local container = k.core.v1.container;
local service = k.core.v1.service;

local redisContainer =
  container.default("redis-master", "redis:2.8.23") +
  container.ports([{name: "redis-server", containerPort: 6379}]);

local redisLabel = {app: "redis", role: "master"};

local redisRC =
  rc.default("redis-master") +
  {spec+: {selector+: redisLabel}} +
  {spec+: {template+: {spec+: {containers: [redisContainer]}}}} +
  {spec+: {template+: {metadata: {labels: redisLabel}}}};

local redisSvc =
  service.default("redis-master") +
  service.mixin.spec.ports([{port: 6379, targetPort: "redis-server"}]) +
  service.mixin.spec.selector(redisLabel) +
  service.mixin.metadata.labels(redisLabel);

{
  redisRC : util.prune(redisRC),
  redisSvc: util.prune(redisSvc),
}
