import std.io
import std.container.string
import std.c

@extern(main)
func main()u8{
    c::setlocale(c::LC_ALL, c"")
    let s = string::new("Hello World\n")
    io::print(s)
    return 0
}