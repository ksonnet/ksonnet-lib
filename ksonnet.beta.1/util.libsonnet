{
  // Remove all empty objects and arrays
  prune(thing)::
    if std.type(thing) == "array" then
      self.pruneArray(thing)
    else if std.type(thing) == "object" then
      self.pruneObj(thing)
    else
      thing,
  // Does this value have real content?
  isContent(v)::
    if v == null then
      false
    else if std.type(v) == "array" then
      std.length(v) > 0
    else if std.type(v) == "object" then
      std.length(v) > 0
    else
      true,
  // Remove all fields that have empty content
  pruneObj(obj):: {
    [x]: $.util.prune(obj[x])
    for x in std.objectFields(obj) if self.isContent(self.prune(obj[x]))
  },
  // Remove all members that have empty content
  pruneArray(arr)::
    [ $.util.prune(x) for x in arr if self.isContent($.util.prune(x)) ]
}