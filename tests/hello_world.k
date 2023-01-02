import std.io
import std.container.string
import std.c

@extern(main)
func main()u8{
    c::setlocale(c::LC_ALL, c"")
    io::println(string::new("Hello World"))
    return 0
}