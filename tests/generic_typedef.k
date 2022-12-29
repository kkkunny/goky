type A<T> T
type B<T> [2]T
type C<T> func(T)T
type D<T> struct{
    data: *B::<T>
}

@extern(main)
func main()u8{
    let a: A::<i32> = 1
    let b: B::<i32> = [1, 2]
    let c: C::<i32> = null
    //let d: D::<A::<i32>> = {}
    return 0
}