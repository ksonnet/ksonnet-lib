local kubeAssert = import "internal/assert.libsonnet";
local base = import "internal/base.libsonnet";
local meta = import "internal/meta.libsonnet";

{
  mixin:: {
    // NOTE: Convenience mixins, will eventually be moved to a mixin
    // package.

    Metadata(mixin):: {metadata+: mixin},
  },

  // A collection of common fields in the Kubernetes API objects,
  // that we do not want to expose for public use. For example,
  // `Kind` appears frequently in API objects of both
  // `extensions/v1beta1` and `v1`, but we don't want users to mess
  // mess with an object's `Kind`.
  local common = {
    Kind(kind):: kubeAssert.Type("kind", kind, "string") {kind: kind},

    // TODO: This sets the metadata property, rather than doing a
    // mixin. Is this what we want?
    Metadata(metadata={}):: {metadata: $.v1.metadata.Default() + metadata},
  },

  v1:: {
    local bases = {
      ConfigMap: base.New("configMap", "AC74E727-0605-4872-8F30-E5CAFB2A0984"),
      Container: base.New("container", "50281784-097C-46A9-8D2C-C6E9078D77D4"),
      ContainerPort:
        base.New("containerPort", "2854EB13-644C-4FEF-A62D-DBAC554D6A24"),
      Metadata: base.New("metadata", "027AE69D-1DD6-42D2-AD47-8F4A55DF9D76"),
      PersistentVolume:
        base.New("persistentVolume", "03113473-7083-4D07-A7FE-83699EB4128C"),
      PersistentVolumeClaim:
        base.New("persistentVolumeClaim", "CD58B997-FF5E-4ED9-8F8A-573E92336D35"),
      Pod: base.New("pod", "2854EB13-644C-4FEF-A62D-DBAC554D6A24"),
      Probe: base.New("probe", "943CF775-B17F-4D25-A794-7D800F08E7FE"),
      Secret: base.New("secret", "0C3D2362-968B-4751-BF67-D58ADA1FC5FC"),
      Service: base.New("service", "87EE499C-EC06-421D-9450-EFE0701851EB"),
      ServicePort: base.New("servicePort", "C38839B7-DA05-4845-B643-E6826E38EA1B"),
      Mount: base.New("mount", "D1E2E601-E64A-4A95-A15C-E78CA724764C"),
      Namespace: base.New("namespace", "6A94A118-F6A7-40EE-8BA1-6096CEC7BDE3"),
    },

    ApiVersion:: { apiVersion: "v1" },

    metadata:: {
      Default(name=null, namespace=null, annotations=null, labels=null)::
        bases.Metadata +
        (if name != null then self.Name(name) else {}) +
        (if namespace != null then self.Namespace(namespace) else {}) + {
          annotations: if annotations == null then {} else annotations,
          labels: if labels == null then {} else labels,
        },

      Name:: CreateNameFunction(false),
      Label:: CreateLabelFunction(false),
      Labels:: CreateLabelsFunction(false),
      Namespace:: CreateNamespaceFunction(false),
      Annotation:: CreateAnnotationFunction(false),
      Annotations:: CreateAnnotationsFunction(false),

      // TODO: Consider renaming this or moving it. `mixins` is
      // probably not something we want to expose to users, at least
      // in this form.
      mixins:: {
        Name:: CreateNameFunction(true),
        Label:: CreateLabelFunction(true),
        Labels:: CreateLabelsFunction(true),
        Namespace:: CreateNamespaceFunction(true),
        Annotation:: CreateAnnotationFunction(true),
        Annotations:: CreateAnnotationsFunction(true),
      },

      //
      // Helpers.
      //

      local CreateNameFunction(isMixin) =
        meta.MixinPartial1(
          isMixin,
          $.mixin.Metadata,
          function(name)
            base.Verify(bases.Metadata) +
            kubeAssert.Type("name", name, "string") +
            {name: name}),

      local CreateLabelFunction(isMixin) =
        meta.MixinPartial2(
          isMixin,
          $.mixin.Metadata,
          function(key, value)
            base.Verify(bases.Metadata) +
            {labels+: {[key]: value}}),

      local CreateLabelsFunction(isMixin) =
        meta.MixinPartial1(
          isMixin,
          $.mixin.Metadata,
          function(labels)
            base.Verify(bases.Metadata) +
            {labels+: labels}),

      local CreateNamespaceFunction(isMixin) =
        meta.MixinPartial1(
          isMixin,
          $.mixin.Metadata,
          function(namespace)
            base.Verify(bases.Metadata) +
            kubeAssert.Type("namespace", namespace, "string") +
            {namespace: namespace}),

      local CreateAnnotationFunction(isMixin) =
        meta.MixinPartial2(
          isMixin,
          $.mixin.Metadata,
          function(key, value)
            base.Verify(bases.Metadata) +
            {annotations+: {[key]: value}}),

      local CreateAnnotationsFunction(isMixin) =
        meta.MixinPartial1(
          isMixin,
          $.mixin.Metadata,
          function(annotations)
            base.Verify(bases.Metadata) +
            {annotations+: annotations}),
    },

    //
    // Namespace.
    //

    namespace:: {
      Default(name)::
        bases.Namespace +
        kubeAssert.Type("name", name, "string") +
        $.v1.ApiVersion +
        common.Kind("Namespace") +
        common.Metadata($.v1.metadata.Name(name)),
    },

    //
    // Ports.
    //
    port:: {
      local protocolOptions = std.set(["TCP", "UDP"]),

      local PortProtocol(protocol, targetBase) =
        kubeAssert.InSet("protocol", protocol, protocolOptions) +
        base.Verify(targetBase) {
          protocol: protocol,
        },

      local PortName(name, targetPort) =
        base.Verify(targetPort) +
        kubeAssert.Type("name", name, "string") {
          name: name,
        },

      container:: {
        Default(containerPort)::
          bases.ContainerPort +
          kubeAssert.ValidPort("containerPort", containerPort) {
            containerPort: containerPort,
          },

        Named(name, containerPort)::
          kubeAssert.Type("name", name, "string") +
          self.Default(containerPort) +
          self.Name(name),

        Name(name):: PortName(name, bases.ContainerPort),

        Protocol(protocol):: PortProtocol(protocol, bases.ContainerPort),

        HostPort(hostPort)::
          base.Verify(bases.ContainerPort) +
          kubeAssert.ValidPort("hostPort", hostPort) {
            hostPort: hostPort
          },

        HostIp(hostIp)::
          base.Verify(bases.ContainerPort) +
          kubeAssert.Type("hostIp", hostIp, "string") {
            hostIP: hostIp,
          },
      },

      service:: {
        Default(servicePort)::
          bases.ServicePort +
          kubeAssert.ValidPort("servicePort", servicePort) {
            port: servicePort,
          },

        WithTarget(servicePort, targetPort)::
          self.Default(servicePort) +
          self.TargetPort(targetPort),

        Named(name, servicePort, targetPort)::
          kubeAssert.Type("name", name, "string") +
          self.Default(servicePort) +
          self.Name(name) +
          self.TargetPort(targetPort),

        Name(name):: PortName(name, bases.ServicePort),

        Protocol(protocol):: PortProtocol(protocol, bases.ServicePort),

        TargetPort(targetPort)::
          base.Verify(bases.ServicePort) {
            // TODO: Assert clusterIP is not set?
            targetPort: targetPort,
          },

        NodePort(nodePort)::
          base.Verify(bases.ServicePort) {
            nodePort: nodePort,
          },
      },
    },

    //
    // Service.
    //
    service:: {
      Default(name, portList, labels={}, annotations={})::
        local defaultMetadata =
          common.Metadata(
            $.v1.metadata.Name(name) +
            $.v1.metadata.Labels(labels) +
            $.v1.metadata.Annotations(annotations));
        local serviceKind = common.Kind("Service");
        bases.Service + $.v1.ApiVersion + serviceKind + defaultMetadata {
          spec: {
            ports: portList,
          },
        },

        Metadata:: $.mixin.Metadata,

        mixin:: {
          metadata:: $.v1.metadata.mixins,

          spec:: {
            Port:: CreatePortFunction(true),
            Selector:: CreateSelectorFunction(true),
            ClusterIp:: CreateClusterIpFunction(true),
            Type:: CreateTypeFunction(true),
            ExternalIps:: CreateExternalIpsFunction(true),
            SessionAffinity:: CreateSessionAffinityFunction(true),
            LoadBalancerIp:: CreateLoadBalancerIpFunction(true),
            LoadBalancerSourceRanges::
              CreateLoadBalancerSourceRangesFunction(true),
            ExternalName:: CreateExternalNameFunction(true),
          },
        },

        //
        // Service spec.
        //

        local typeOptions = std.set([
          "ExternalName", "ClusterIP", "NodePort", "LoadBalancer"]),
        local sessionAffinityOptions = std.set(["ClientIP", "None"]),

        spec:: {
          Port:: CreatePortFunction(false),
          Selector:: CreateSelectorFunction(false),
          ClusterIp:: CreateClusterIpFunction(false),
          Type:: CreateTypeFunction(false),
          ExternalIps:: CreateExternalIpsFunction(false),
          SessionAffinity:: CreateSessionAffinityFunction(false),
          LoadBalancerIp:: CreateLoadBalancerIpFunction(false),
          LoadBalancerSourceRanges::
            CreateLoadBalancerSourceRangesFunction(false),
          ExternalName:: CreateExternalNameFunction(false),
        },

        local specMixin(mixin) = { spec+: mixin },

        local CreatePortFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(port)
              // base.Verify(bases.Service) +
              {ports+: [port]}),

        local CreateSelectorFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(selector)
              // base.Verify(bases.Service) +
              {selector: selector}),

        local CreateClusterIpFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(clusterIp)
              // base.Verify(bases.Service) +
              kubeAssert.Type("clusterIp", clusterIp, "string") +
              {clusterIP: clusterIp}),

        local CreateTypeFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(type)
              // base.Verify(bases.Service) +
              kubeAssert.InSet("type", type, typeOptions) +
              {type: type}),

        local CreateExternalIpsFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(externalIpList)
              // base.Verify(bases.Service) +
              // TODO: Verify that externalIpList is a list of string.
              kubeAssert.Type("externalIpList", externalIpList, "array") +
              {externalIPs: externalIpList}),

        local CreateSessionAffinityFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(sessionAffinity)
              // base.Verify(bases.Service) +
              kubeAssert.InSet(
                "sessionAffinity", sessionAffinity, sessionAffinityOptions) +
              {sessionAffinity: sessionAffinity}),

        local CreateLoadBalancerIpFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(loadBalancerIp)
              // base.Verify(bases.Service) +
              kubeAssert.Type("loadBalancerIp", loadBalancerIp, "string") +
              {loadBalancerIP: loadBalancerIp}),

        local CreateLoadBalancerSourceRangesFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(loadBalancerSourceRanges)
              // base.Verify(bases.Service) +
              // TODO: Verify that loadBalancerSourceRanges is a list
              // of string.
              kubeAssert.Type(
                "loadBalancerSourceRanges", loadBalancerSourceRanges, "array") +
              {loadBalancerSourceRanges: loadBalancerSourceRanges}),

        local CreateExternalNameFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(externalName)
              // base.Verify(bases.Service) +
              kubeAssert.Type("externalName", externalName, "string") +
              {externalName: externalName}),
    },

    configMap:: {
      Default(namespace, configMapName, data):
        bases.ConfigMap +
        $.v1.ApiVersion +
        common.Kind("ConfigMap") +
        common.Metadata(
          $.v1.metadata.Name(configMapName) +
          $.v1.metadata.Namespace(namespace)) {
          data: data,
        },

      DefaultFromClaim(namespace, name, claim)::
        self.Default(namespace, name, claim.metadata.name)
    },

    secret:: {
      Default(namespace, configMapName, data)::
        bases.Secret +
        $.v1.ApiVersion +
        common.Kind("Secret") +
        common.Metadata(
          $.v1.metadata.Name(configMapName) +
          $.v1.metadata.Namespace(namespace)) {
          data: data,
        },

      StringData(stringData)::
        base.Verify(bases.Secret) {
          stringData: stringData,
        },

      Type(type)::
        base.Verify(bases.Secret) +
        kubeAssert.Type("type", type, "string") {
          type: type,
        },
    },

    //
    // Volume.
    //

    //
    // NOTE: TODO: YOU ARE HERE. You haven't implemented type checking
    // beyond this point.
    //

    volume:: {
      persistent:: {
        // TODO: Add checks to the parameters here.
        Default(name, claimName):: bases.PersistentVolume {
          name: name,
          persistentVolumeClaim: {
            claimName: claimName,
          },
        },

        DefaultFromClaim(name, claim)::
          self.Default(name, claim.metadata.name)
      },

      hostPath:: {
        // TODO: Add checks to the parameters here.
        Default(name, path):: {
          name: name,
          hostPath: {
            path: path
          },
        },
      },

      configMap:: {
        // TODO: Add checks to the parameters here.
        Default(name, configMapName):: {
          name: name,
          configMap: {
            name: configMapName,
          },
        },
      },

      secret:: {
        // TODO: Add checks to the parameters here.
        Default(name, secretName):: {
          name: name,
          secret: {
            secretName: secretName,
          },
        },
      },

      emptyDir:: {
        Default(name):: {
          name: name,
          emptyDir: {},
        },
      },

      //
      // Mount.
      //
      mount:: {
        Default(name, mountPath, readOnly=false):: bases.Mount {
          name: name,
          mountPath: mountPath,
          readOnly: readOnly,
        },

        FromVolume(volume, mountPath, readOnly=false)::
          self.Default(volume.name, mountPath, readOnly),

        FromConfigMap(configMap, mountPath, readOnly=false)::
          self.Default(configMap.name, mountPath, readOnly),
      },

      //
      // Claim.
      //
      claim:: {
        // TODO: This defaults to a storage class; probably it
        // shouldn't.
        DefaultPersistent(
          namespace,
          claimName,
          accessModes,
          size,
          storageClass="fast"
        ):
          local defaultMetadata =
            common.Metadata(
              $.v1.metadata.Name(claimName) +
              $.v1.metadata.Namespace(namespace) +
              $.v1.metadata.Annotations({
                "volume.beta.kubernetes.io/storage-class": storageClass,
              }));
          bases.PersistentVolumeClaim +
          $.v1.ApiVersion +
          common.Kind("PersistentVolumeClaim") +
          defaultMetadata {
            // TODO: Move this assert to `kubeAssert.Type`.
            assert std.type(accessModes) == "array"
              : "'accessModes' must by of type 'array'",
            spec: {
              accessModes: accessModes,
              resources: {
                requests: {
                  storage: size
                },
              },
            },
          },
      },
    },

    //
    // Probe.
    //
    probe:: {
      local defaultTimeout = 1,
      local defaultPeriod = 10,

      Default(
        initDelaySecs,
        timeoutSecs=defaultTimeout,
        periodSeconds=defaultPeriod
      ):: bases.Probe {
        initialDelaySeconds: initDelaySecs,
        timeoutSeconds: timeoutSecs,
      },

      Http(
        getPath,
        portName,
        initDelaySecs,
        timeoutSecs=defaultTimeout,
        periodSeconds=defaultPeriod
      ):: self.Default(initDelaySecs, timeoutSecs) {
          httpGet: {
            path: getPath,
            port: portName,
          },
        },

      Tcp(
        port,
        initDelaySecs,
        timeoutSecs=defaultTimeout,
        periodSeconds=defaultPeriod
      ):: self.Default(initDelaySecs, timeoutSecs) {
          tcpSocket: {
            port: port,
          },
        },

      Exec(
        command,
        initDelaySecs,
        timeoutSecs=defaultTimeout,
        periodSeconds=defaultPeriod
      ):: self.Default(initDelaySecs, timeoutSecs) {
          exec: {
            command: command,
          },
        },
    },

    //
    // Container.
    //
    container:: {
      local imagePullPolicyOptions = std.set(["Always", "Never", "IfNotPresent"]),

      Default(name, image, imagePullPolicy="Always")::
        bases.Container +
        // TODO: Make "Always" the default only when we're doing the :latest.
        kubeAssert.Type("name", name, "string") +
        kubeAssert.Type("image", image, "string") +
        kubeAssert.InSet("imagePullPolicy", imagePullPolicy, imagePullPolicyOptions) {
          name: name,
          image: image,
          imagePullPolicy: imagePullPolicy,
          // TODO: Think carefully about whether we want an empty list here.
          ports: [],
          env: [],
        },

      Command(command):: base.Verify(bases.Container) {
        command: command,
      },

      // TODO: Should this take a k/v pair instead?
      Env(env):: base.Verify(bases.Container) {
        env: env,
      },

      Resources(resources):: base.Verify(bases.Container) {
        resources: resources
      },

      Ports(ports):: base.Verify(bases.Container) {
        ports: ports,
      },

      Port(port):: base.Verify(bases.Container) { ports+: [port] },

      NamedPort(name, port):: base.Verify(bases.Container) {
        ports+: [$.v1.port.container.Named(name, port)],
      },

      LivenessProbe(probe):: base.Verify(bases.Container) {
        livenessProbe: probe,
      },

      ReadinessProbe(probe):: base.Verify(bases.Container) {
        readinessProbe: probe,
      },

      VolumeMounts(mounts):: base.Verify(bases.Container) {
        volumeMounts: mounts,
      },
    },

    //
    // Env.
    //
    env:: {
      Variable(name, value):: {
        name: name,
        value: value,
      },

      ValueFrom(name, configMapName, configMapKey):: {
        name: name,
        valueFrom: {
          configMapKeyRef: {
            name: configMapName,
            key: configMapKey,
          },
        },
      },

      ValueFromSecret(name, secretName, secretKey):: {
        name: name,
        valueFrom: {
          secretKeyRef: {
            name: secretName,
            key: secretKey,
          },
        },
      },
    },

    //
    // Pods.
    //
    pod:: {
      Default(spec)::
        bases.Pod +
        $.v1.ApiVersion +
        common.Kind("Pod") +
        common.Metadata() {
          spec: spec,
        },

      Metadata:: $.mixin.Metadata,

      // TODO: Consider making this just a function on the pod itself.
      template:: {
        // TODO: This does not really belong here. We should have
        // something like `deployment.spec.Template` instead.
        Default(spec)::
          common.Metadata() {
            spec: spec,
          },

        Metadata:: $.mixin.Metadata,
        mixin:: {metadata:: $.v1.metadata.mixins},
      },

      // TODO: Consider making this just a function on the pod itself.
      // TODO: Shouldn't this just be in pod.template?
      spec:: {
        Containers(containers):: {
          containers: containers,
        },

        Volumes(volumes):: {
          volumes: volumes,
        },

        DnsPolicy(policy="ClusterFirst"):: {
          dnsPolicy: policy,
        },

        RestartPolicy(policy="Always"):: {
          restartPolicy: policy,
        },
      },
    },
  },

  extensions:: {
    v1beta1: {
      local bases = {
        Deployment: base.New("deployment", "176A7BEF-E577-4EBD-952D-5E8F7BB7AE1A"),
      },

      ApiVersion:: { apiVersion: "extensions/v1beta1" },

      //
      // Deployments.
      //
      deployment:: {
        Default(name, spec)::
          local defaultMetadata =
            common.Metadata($.v1.metadata.Name(name));
          bases.Deployment +
          $.extensions.v1beta1.ApiVersion +
          common.Kind("Deployment") +
          defaultMetadata {
            spec: spec,
          },

        Metadata:: $.mixin.Metadata,

        mixin:: {
          metadata:: $.v1.metadata.mixins,

          podTemplate:: {
            local templateMixin(mixin) = {
              // TODO: Add base verification here.
              spec+: {
                template+: {
                  spec+: mixin
                },
              },
            },

            Volumes(volumes)::
              templateMixin($.v1.pod.spec.Volumes(volumes)),

            Containers(containers)::
              templateMixin($.v1.pod.spec.Containers(containers)),

            // TODO: Consider moving this default to some common
            // place, so it's not duplicated.
            DnsPolicy(policy="ClusterFirst")::
              templateMixin($.v1.pod.spec.DnsPolicy(policy=policy)),

            RestartPolicy(policy="Always")::
              templateMixin($.v1.pod.spec.RestartPolicy(policy=policy)),
          },

          spec:: {
            Selector:: CreateSelectorFunction(true),
            MinReadySeconds:: CreateMinReadySecondsFunction(true),
            RollingUpdateStrategy::
              CreateRollingUpdateStrategyFunction(true),
          },
        },

        // TODO: Consolidate so that we have one of these functions.
        local specMixin(mixin) = {spec+: mixin},

        NodeSelector(labels)::
          base.Verify(bases.Deployment) +
          specMixin({nodeSelector: labels}),

        spec:: {
          ReplicatedPod(replicas, podTemplate):: {
            replicas: replicas,
            template: podTemplate,
          },

          Selector:: CreateSelectorFunction(false),
          MinReadySeconds:: CreateMinReadySecondsFunction(false),
          RollingUpdateStrategy::
            CreateRollingUpdateStrategyFunction(false),
        },

        local CreateSelectorFunction(isMixin) =
          meta.MixinPartial1(
            isMixin,
            specMixin,
            function(labels)
              // base.Verify(bases.Service) +
              {
                selector: {
                  matchLabels: labels,
                },
              }),

        local CreateMinReadySecondsFunction(isMixin) =
          local partial =
            meta.MixinPartial1(
              isMixin,
              specMixin,
              function(seconds)
                // base.Verify(bases.Service) +
                {minReadySeconds: seconds});
          function(seconds=0) partial(seconds),

        local CreateRollingUpdateStrategyFunction(isMixin) =
          local partial =
            meta.MixinPartial2(
              isMixin,
              specMixin,
              // function(maxSurge=1, maxUnavailable=1)
              function(maxSurge, maxUnavailable)
                // base.Verify(bases.Service)
                {
                  strategy: {
                    rollingUpdate: {
                      maxSurge: maxSurge,
                      maxUnavailable: maxUnavailable,
                    },
                    type: "RollingUpdate",
                  },
                });
          function(maxSurge=1, maxUnavailable=1)
            partial(maxSurge, maxUnavailable),
      },

      IngressSpec(domain, serviceName, servicePort):: {
        rules: [
          {
            host: domain,
            http: {
              paths: [{
                backend: {
                  serviceName: serviceName,
                  servicePort: servicePort,
                }}]
            }
          }
        ]
      },
    },
  },
}