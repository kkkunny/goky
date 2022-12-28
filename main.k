type ArrayList<T> struct{
    cap: usize
    len: usize
    ptr: *T
}

func add<T>(x: T, y: T)T{
    return x + y
}

@extern(main)
func main()u8{
    let a: ArrayList::<i32> = {}
    return add::<u8>(1, 2)
}
