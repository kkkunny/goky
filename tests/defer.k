import std.c

@extern(main)
func main()u8{
    defer c::puts(c"exit")
    return 0
}