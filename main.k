@extern(printf)
func printf(*i8)

func print(s: [13]i8){
    printf(&s as *i8)
}

@extern(main)
func main()i32{
    print(c"hello world\n")
    return 0
}
