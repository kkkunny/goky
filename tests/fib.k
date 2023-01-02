import std.c

func fib(n: i32)i32{
    if n == 1 || n == 2{
        return 1
    }
    return fib(n-1) + fib(n-2)
}

@extern(main)
func main()u8{
    let fileStr = c"fib.k"
    let funcStr = c"main"
    let str = c"fib(5) != 5"
    if fib(5) != 5{
        c::__assert_fail(&str as *c::char, &fileStr as *c::char, 16, &funcStr as *c::char)
    }
    return 0
}