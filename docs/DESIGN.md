# Design

> διὸ δὴ συμμειγνύμενα αὐτά τε πρὸς αὑτὰ καὶ πρὸς ἄλληλα τὴν ποικιλίαν ἐστὶν ἄπειρα: ἧς δὴ δεῖ θεωροὺς γίγνεσθαι τοὺς μέλλοντας περὶ φύσεως εἰκότι λόγῳ χρήσεσθαι.

> *So their combinations with themselves and with each other give rise
> to endless complexities, which anyone who is to give a likely
> account of reality must survey.*

Plato, in [*Timaeus*][timaeus], speaking about the Platonic solids,
which were then viewed as the idealized primary constituents of the
physical universe.

## Goals

More particularly, we would like to:

* Allow users to template common patterns, without making it hard to
  customize API objects when they differ from the pattern.
* Make it harder to declare an invalid configuration.
* Create abstractions that are amenable to tooling (_e.g._, syntax
  highlighting, static analysis, autocompletion, and so on).
* Allow users tools to clearly separate data declaration from
  application configuration (_e.g._ through the use of simple
  variables, _etc_.).

Non-goals:

* Make a higher-level API than the Kubernetes API. We believe that it
  will become much easier to develop higher-level APIs if we develop
  good tooling for managing the Kubernetes API that already exists.



This pattern of using `+` to incrementally build up Kubernetes API
objects is very common, and used throughout `ksonnet-lib`. This
pattern has several distinct advantages:

* **Makes it very likely that you will generate a correct
  specification.** `container.Image` takes a name for the container
  and a name for the image, and spits out a well-formed container
  spec. Your only opportunity to mess this up is to put the wrong data
  in!
* **High fidelity translation to the underlying Kubernetes object
  API.** Because we are building up real Kubernetes API objects
  directly, you have the full expressive power of the Kubernetes APIs,
  not a lossy higher-level API. Because of this, you can template what
  you want, and then fall back on normal JSON when you don't, because
  there's no magical translation layer.
* **Safety.** Each of these functions is able to check that its parent
  is spec-compliant, so you won't have to worry about `+`'ing things
  together that don't make sense.
* **Tooling support and static analysis (eventually).** For example,
  you will be able to dot into `container` to see what the available
  fields are.



[timaeus]: http://www.perseus.tufts.edu/hopper/text?doc=Perseus%3Atext%3A1999.01.0179%3Atext%3DTim.%3Asection%3D57d
