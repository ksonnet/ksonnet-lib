# 1: Jsonnet recap: references, variables, simple JSON templating

Before we demonstrate the core abstractions of `kube.libsonnet`, it is
worth spending time on a whirlwind tour of Jsonnet, to familiarize
ourselves with the features we will use in this tutorial. (_If you
already know Jsonnet, you will lose nothing by skipping to the [next
section][k8s-prims]._)

For the purposes of this tutorial, you can think of Jsonnet as a
domain-specific language meant to make it easy to declare and template
languages. Think JSON, but with:

* variables (both [lexically-scoped locals][jsonnetLocals] and
  JsonPath-style [references][jsonnetReferences])
* [functions][jsonnetFunctions]
* the ability to define libraries and [import][jsonnetImports] them
* some notion of [object-oriented inheritance between JSON
  objects][jsonnetOO]
* and a bunch of the [syntax annoyances ironed out][jsonnetSyntax].

For the purposes of this tutorial, you need only a very small subset
of these concepts. They are:

### Local variables and references

In Jsonnet, it is possible to define lexically-scoped local variables:

```c++
{
  local foo = "bar",
  baz: foo,
}
```

which produces:

```json
{ "baz": "bar" }
```

Jsonnet additionally exposes a `self` to access properties of the
current object, and a JsonPath-style `$`, which refers to the "root
object" (or: the grandparent who is farthest away from the `$`):

```c++
{
  foo: "bar",
  baz: self.foo,
  cow: {
    moo: $.foo,
  },
}
```

```json
{
  "foo": "bar",
  "baz": "bar",
  "cow": { "moo": "bar" }
}
```

It is worth noting that both `local` variables and references are
_order-independent_, which is a decision that largely falls out of
JSON's design. Notice, for example, that if we re-order `foo` and
`baz`, it does not affect the output of Jsonnet:

```c++
{
  baz: self.foo,
  cow: {
    moo: $.foo,
  },

  // This is perfectly legal.
  foo: "bar",
}
```

### Functions

Jsonnet implements lexically-scoped functions, but they can be
declared in a few ways, and it's worth pointing them out.

In the example below, note the use of the double colon (`::`) in
the declaration of `function2`. This marks the field as _hidden_,
which is a concept we will look closer at in the section on
object-orientation. For now, it is only important to understand that a
function must be either `local` or hidden with `::`, because Jsonnet
doesn't know how to render a function as JSON data. (Instead of
rendering it, Jsonnet will complain and crash.)

```c++
{
  local function1(arg1) = { foo: arg1 },
  function2(arg1="cluck"):: { bar: arg1 },
  cow: function1("moo"),
  chicken: self.function2(),
}
```

```json
{
   "chicken": {
      "bar": "cluck"
   },
   "cow": {
      "foo": "moo"
   }
}
```

### Object-orientation (inheritance, mixins, _etc_.)

One of Jsonnet's most powerful features, which we use liberally in
this tutorial and in `kube.libsonnet`, is its object model, which
implements a concise, [well-specified _algebra_][jsonnetAlgebra] for
combining JSON-like objects.

The primary tool for combining objects is the `+` operator. In this
example we see two objects (the first is called the _parent_, or
_base_, and the second is called the _child_) that are combined with
the `+`. The child (which is said to _inherit_ from the parent)
overwrites the `bar` property that was defined in the parent:

```c++

{
  // Parent object.
  foo: "foo",
  bar: "bar",
} + {
  // Child object.
  bar: "fubar",
}
```

```json
{
   "bar": "fubar",
   "foo": "foo"
}
```

It is sometimes convenient for a child to reference members of the
parent, so Jsonnet also exposes `super`, which behaves a lot like
`self`, except in reference to the parent:

```c++
{
  foo: "foo",
} + {
  bar: super.foo + "bar",
}
```

```json
{
   "bar": "foobar",
   "foo": "foo"
}
```

One interesting aspect of `super` is that it can be "mixed in",
meaning that if you have an object that refers to `super.bar`, then it
can dynamically be made to inherit from _any object_ that has a `bar`
property. For example:

```c++
local fooTheBar = { bar: super.bar + "foo" };
{
  bar: "bar",
} + fooTheBar
```

```json
{
   "bar": "barfoo"
}
```

This stands in contrast to the object model of (say) Java, where you
would have to declare at compile time an `Animal` class before a `Dog`
class could be made to inherit from it. The technique above (called a
_mixin_) causes the object to inherit dynamically, at runtime rather
than compile time.

Lastly, Jsonnet allows you to create hidden properties, not included
when we generate the final JSON. Denoted with with a `::`, they are
also visible to all descendent objects (_i.e._, children,
grandchildren, _etc_.), and are useful for holding data you'd like to
use to construct other properties, but not expose as part of the
generated JSON itself:

```c++
{
  foo:: "foo",
} + {
  bar: super.foo + "bar",
}
```

```json
{
   "bar": "foobar"
}
```
