import std.os
import std.c

@extern(main)
func main()u8{
    let s = c"hello world"
    c:__c_puts(&s as *i8)
    return 0
}
