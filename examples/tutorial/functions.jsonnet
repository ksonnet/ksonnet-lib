{
  local function1(arg1) = { foo: arg1 },
  function2(arg1="cluck"):: { bar: arg1 },
  cow: function1("moo"),
  chicken: self.function2(),
}