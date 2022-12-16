import std.sys
import std.fmt

@extern(_start)
func main(){
    let s = c"hello world\n"
    fmt:print(&s as *i8)
    sys:exit(12)
}
