@extern(printf)
func printf(*i8)

@extern(main)
func main()i32{
    let s = c"hello world\n";
    defer printf(&s as *i8);
    return 0
}
