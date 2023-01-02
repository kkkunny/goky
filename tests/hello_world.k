import std.c

@extern(main)
func main()u8{
    let fileStr = c"hello_world.k"
    let funcStr = c"main"
    let s = c"hello world"
    let str = c"c::puts(&s as *c::char) != 12"
    if c::puts(&s as *c::char) != 12{
        c::__assert_fail(&str as *c::char, &fileStr as *c::char, 10, &funcStr as *c::char)
    }
    return 0
}