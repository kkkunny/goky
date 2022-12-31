@inline(true)
func add(x: i32, y: i32)i32{
    return x + y
}

@extern(main)
func main()u8{
    return add(1, 2) as u8
}