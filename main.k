type A bool

@extern(main)
func main()u8{
    let a: A = true
    if a{
        return 1
    }
    return 0
}
