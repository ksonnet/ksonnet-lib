{
  MixinPartial1(isMixin, createMixin, fn)::
    if isMixin
    then function(arg1) createMixin(fn(arg1))
    else fn,

  MixinPartial2(isMixin, createMixin, fn)::
    if isMixin
    then function(arg1, arg2) createMixin(fn(arg1, arg2))
    else fn,
}