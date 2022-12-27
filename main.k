import std.c

@extern(main)
func main()u8{
    let i: c:size_t = 10
    return i as u8
}
