local domain = '.example.org';
local services = ['one', 'two'];
local computedServices = [
  local serviceName = service;
  local serviceUrl = service + domain;

  { [serviceName]: serviceUrl }
  for service in services
];

std.prune(computedServices)