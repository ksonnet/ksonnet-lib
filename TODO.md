# TODO list

* [ ] Decide how we want `service.Port` to work with `port.service.*`.
  Perhaps we should import all those functions into this namespace,
  for convenience sake? Same for ports, volume mounts, etc.
* [ ] Rethink whether `container.Env` should take a whole list, or a k/v pair.
* [ ] Consider erroring out if a name is not DNS conformant. (e.g., in
  a Deployment).
