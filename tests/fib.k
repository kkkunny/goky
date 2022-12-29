func fib(n: i32)i32{
    if n == 1 || n == 2{
        return 1
    }
    return fib(n-1) + fib(n-2)
}

@extern(main)
func main()u8{
    let n = fib(5)
    return 0
}