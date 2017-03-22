# TODO list

* [ ] Decide how we want `service.Port` to work with `port.service.*`.
  Perhaps we should import all those functions into this namespace,
  for convenience sake? Otherwise we end up repeating ourselves a lot.
  Same for ports, volume mounts, etc.
  * [ ] Consider moving these helpers into util.
* [ ] Rethink whether `container.Env` should take a whole list, or a k/v pair.
* [ ] Consider erroring out if a name is not DNS conformant. (e.g., in
  a Deployment).
* [ ] Rethink the namespace issues. Importing util should probably
  merge itself with the core namespace?
* [ ] Figure out pod.template namespace. If we merge the pod
  namespaces, we lose the pod.template namespace.
* [ ] Unify different kinds of volumes. No need for a persistent
  volume with its own base (for example).
* [ ] Get rid of default storage class in PVC.
* [ ] metadata.Labels currently overwrites whatever came before it. Desired?
