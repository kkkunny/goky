import std.c

type A i32

func (A) get()i32{
    return *self as i32
}

@extern(main)
func main()u8{
    let a: A = 1
    if a.get() != 1{
        c::__assert_fail(c"a.get() != 1", c"method.k", 13, c"main")
    }
    return 0
}