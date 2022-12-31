type A i32

func (A) get()i32{
    return *self as i32
}

@extern(main)
func main()u8{
    let a: A = 1
    return a.get() as u8
}