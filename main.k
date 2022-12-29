import std.c

@extern(main)
func main()u8{
    let a: c:int = 1
    return a as u8
}
