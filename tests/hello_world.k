import std.c

@extern(main)
func main()u8{
    if c::puts(c"hello world") != 12{
        c::__assert_fail(c"c::puts(c\"hello world\") != 12", c"hello_world.k", 6, c"main")
    }
    return 0
}