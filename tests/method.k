import std.c

type A i32

func (A) get()i32{
    return *self as i32
}

@extern(main)
func main()u8{
    let fileStr = c"method.k"
    let funcStr = c"main"
    let a: A = 1
    let str = c"a.get() != 1"
    if a.get() != 1{
        c::__assert_fail(&str as *c::char, &fileStr as *c::char, 16, &funcStr as *c::char)
    }
    return 0
}