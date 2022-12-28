import std.c

@extern(main)
func main()u8{
    let s = c""
    c:setlocale(c:LC_ALL, &s as *c:char)
    c:putwchar('你' as c:wchar_t)
    c:putwchar('好' as c:wchar_t)
    c:putwchar('世' as c:wchar_t)
    c:putwchar('界' as c:wchar_t)
    c:putwchar('\n' as c:wchar_t)
    return 0
}
