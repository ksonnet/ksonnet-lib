{
  global: {
    restart: false,
  },
  components: {
    'guestbook-ui': {
      containerPort: 80,
      image: 'gcr.io/heptio-images/ks-guestbook-demo:0.2',
      name: 'guestbook-ui',
      replicas: 5,
      servicePort: 80,
      type: 'NodePort',
    },
  },
}
