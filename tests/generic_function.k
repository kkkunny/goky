func add<T>(x: T, y: T)T{
    return x + y
}

@extern(main)
func main()u8{
    let a = add::<i32>(1, 2)
    let b = add::<f32>(1, 2)
    return 0
}