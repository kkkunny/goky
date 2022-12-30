import std.c

@extern(main)
func main()u8{
    let s = c"hello world"
    c::puts(&s as *c::char)
    return 0
}