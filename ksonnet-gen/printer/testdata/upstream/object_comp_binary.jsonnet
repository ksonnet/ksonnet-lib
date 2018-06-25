local o = {
  a: 'a',
  b: 'b',
};

{
  ['pre-' + key]: o[key]
  for key in std.objectFields(o)
}