func add[T](x: T, y: T)T{
    return x + y
}

@extern(main)
func main()u8{
    return add::<u8>(1, 2)
}
