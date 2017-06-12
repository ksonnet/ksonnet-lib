local k = import "../../ksonnet.beta.1/k.libsonnet";
local util = import "../../ksonnet.beta.1/util.libsonnet";

local rc = k.core.v1.replicationController;
local container = k.core.v1.container;
local service = k.core.v1.service;

local guestbookContainer =
  container.default("guestbook", "gcr.io/google_containers/guestbook:v3") +
  container.ports([{name: "http-server", containerPort: 3000}]);

local guestbookLabel = {app: "guestbook"};

local guestbookRC =
  rc.default("guestbook") +
  rc.mixin.spec.replicas(3) +
  {spec+: {selector+: guestbookLabel}} +
  {spec+: {template+: {spec+: {containers: [guestbookContainer]}}}} +
  {spec+: {template+: {metadata: {labels: guestbookLabel}}}};

local guestbookSvc =
  service.default("guestbook") +
  service.mixin.spec.ports([{port: 3000, targetPort: "http-server"}]) +
  service.mixin.spec.selector(guestbookLabel) +
  service.mixin.metadata.labels(guestbookLabel) +
  {spec+: {type: "LoadBalancer"}};

{
  guestbookRC : util.prune(guestbookRC),
  guestbookSvc: util.prune(guestbookSvc),
}
