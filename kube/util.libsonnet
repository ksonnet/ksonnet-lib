local core = import "./core.libsonnet";

{
  app:: {
    v1:: {
      container:: {
        NewWithPorts(name, image, ports)::
          core.v1.container.Default(name, image) +
          core.v1.container.Ports(ports),
      },

      env:: {
        array:: {
          // TODO: In all of these, check that we're not duplicating
          // the variables, as the order is independent in Jsonnet,
          // and we will mess it up.

          FromConfigMap(configMap, envSpec)::
            self.FromConfigMapName(configMap.metadata.name, envSpec),

          FromConfigMapName(configMapName, envSpec)::
            [core.v1.env.ValueFrom(name, configMapName, envSpec[name])
              for name in std.objectFields(envSpec)],

          FromSecret(secret, envSpec)::
            self.FromSecretName(secret.metadata.name, envSpec),

          FromSecretName(secretName, envSpec)::
            [core.v1.env.ValueFromSecret(name, secretName, envSpec[name])
              for name in std.objectFields(envSpec)],

          FromObj(envVariables)::
            [core.v1.env.Variable(name, envVariables[name])
              for name in std.objectFields(envVariables)],
        },
      },

      pod:: {
        FromContainer(container, labels={app: container.name})::
          core.v1.pod.Default(core.v1.pod.spec.Containers([container])) +
          core.v1.pod.Labels(labels),

        local mixinSpec(mixin) = {
          spec+: mixin,
        },

        Volumes(volumes):: mixinSpec({volumes: volumes}),

        template:: {
          FromContainer(container, labels={app: container.name}, volumes=[])::
            local spec =
              core.v1.pod.spec.Volumes(volumes) +
              core.v1.pod.spec.Containers([container]);
            core.v1.pod.template.Default(spec) +
            core.v1.pod.template.Labels(labels),
        },
      },

      port:: {
        service:: {
          array:: {
            FromContainerPorts(createServicePort, containerPorts):: [
              core.v1.port.service.Named(
                port.name, createServicePort(port), port.name)
              for port in containerPorts],
          }
        },
      }
    },

    v1beta1:: {
      deployment:: {
        FromPodTemplate(name, replicas, podTemplate, labels={})::
          core.extensions.v1beta1.deployment.Default(
            name,
            core.extensions.v1beta1.deployment.spec.ReplicatedPod(replicas, podTemplate)) +
          core.extensions.v1beta1.deployment.Labels(labels),

        FromContainer(
          name,
          replicas,
          container,
          labels={},
          podLabels={app: container.name},
          volumes=[]
        )::
          self.FromPodTemplate(
            name,
            replicas,
            $.app.v1.pod.template.FromContainer(
              container, labels=podLabels, volumes=volumes),
            labels=labels),

        MixinSpec(spec):: {
          spec+: spec
        },
      },
    },
  },
}