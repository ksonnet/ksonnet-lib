local k8s = import "k8s.libsonnet";

local extensions = k8s.extensions;
local core = k8s.core;

local daemonSet = extensions.v1beta1.daemonSet;
local deployment = extensions.v1beta1.deployment;
local container = deployment.mixin.spec.template.spec.containersType;
local volume = deployment.mixin.spec.template.spec.volumesType;

local hidden = {
  container:: container + {
    new(name, image)::
      super.name(name) +
      super.image(image),

    envType:: container.envType + {
      new(key, value):: super.name(key) + super.value(value),
      fromSecretRef(name, secretRefName, secretRefKey)::
        super.name(name) +
        super.mixin.valueFrom.secretKeyRef.name(secretRefName) +
        super.mixin.valueFrom.secretKeyRef.key(secretRefKey),

      fromFieldPath(name, fieldPath)::
        container.envType.name(name) +
        container.envType.mixin.valueFrom.fieldRef.fieldPath(fieldPath),
    },

    volumeMountsType:: container.volumeMountsType + {
      new(name, mountPath, readonly=false)::
        super.new() +
        super.name(name) +
        super.mountPath(mountPath) +
        super.readOnly(readonly),
    },

    portsType:: container.portsType + {
      named(name, containerPort)::
        super.new() +
        super.name(name) +
        super.containerPort(containerPort),
    }
  },

  volume:: volume + {
    fromPvc(name, claimName)::
      super.new() +
      super.name(name) + {
        persistentVolumeClaim: claimName
      },

    fromHostPath(name, hostPath)::
      volume.name(name) +
      volume.mixin.hostPath.path(hostPath),

    fromConfigMap(name, configMapName=null, items=null)::
      local configMap = volume.mixin.configMap;
      volume.name(name) +
        (if configMapName != null then configMap.name(configMapName) else {}) +
        (if items != null then configMap.items(items) else {}),

    mixin:: volume.mixin + {
      configMap:: volume.mixin.configMap + {
        itemsType:: volume.mixin.configMap.itemsType + {
          new(key, path)::
            super.key(key) +
            super.path(path),
        },
      },
    },
  },
};

k8s + {
  core:: core + {
    v1:: core.v1 + {
      list:: {
        new(items)::
          {apiVersion: "v1"} +
          {kind: "List"} +
          if std.type(items) == "array" then {items+: items} else {items: [items]},
      },

      namespace:: core.v1.namespace + {
        new(name)::
          super.new() +
          super.mixin.metadata.name(name) +
          super.mixin.metadata.labels({"name": name}),
      },

      service:: core.v1.service + {
        new(name, selectorLabels, ports)::
          super.new() +
          super.mixin.metadata.name(name) +
          super.mixin.spec.selector(selectorLabels) +
          super.mixin.spec.ports(ports),

        mixin:: core.v1.service.mixin + {
          spec:: core.v1.service.mixin.spec + {
            portsType:: core.v1.service.mixin.spec.portsType + {
              tcp(servicePort, targetPort)::
                super.new() +
                super.port(servicePort) + {
                  targetPort: targetPort,
                },
            },
          },
        },
      },
    },
  },

  extensions:: extensions + {
    v1beta1:: extensions.v1beta1 + {
      daemonSet:: daemonSet + {

      mapContainers(f):: {
        local podContainers = super.spec.template.spec.containers,
        spec+: {
          template+: {
            spec+: {
              // IMPORTANT: This overwrites the `containers` field
              // for this deployment.
              containers: std.map(f, podContainers),
            },
          },
        },
      },

        mixin:: daemonSet.mixin + {
          spec:: daemonSet.mixin.spec + {
            template:: daemonSet.mixin.spec.template + {
              spec:: daemonSet.mixin.spec.template.spec + {
                containersType:: hidden.container,
                volumesType:: hidden.volume,
              },
            },
          },
        },
      },

      deployment:: extensions.v1beta1.deployment + {
        new(name, replicas, containers, podLabels={})::
          super.new() +
          super.mixin.metadata.name(name) +
          super.mixin.spec.replicas(replicas) +
          super.mixin.spec.template.spec.containers(containers) +
          super.mixin.spec.template.metadata.labels(podLabels),

        mapContainers(f):: {
          local podContainers = super.spec.template.spec.containers,
          spec+: {
            template+: {
              spec+: {
                // IMPORTANT: This overwrites the `containers` field
                // for this deployment.
                containers: std.map(f, podContainers),
              },
            },
          },
        },

        mixin:: deployment.mixin + {
          // extensions.v1beta1.deployment.mixin.spec.template.spec.containersType
          spec:: deployment.mixin.spec + {
            template:: deployment.mixin.spec.template + {
              spec:: deployment.mixin.spec.template.spec + {
                containersType:: hidden.container,
                volumesType:: hidden.volume,
              },
            },
          },
        },
      },
    },
  },
}
