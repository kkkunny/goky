type A struct{
    v: u8
}

@extern(main)
func main()u8{
    let a: A = {1}
    let b = &a
    return b.v
}
