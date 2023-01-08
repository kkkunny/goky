import std.c

@extern(main)
func main()u8{
    let i = "你好世界！"
    let j: [2]i32
    if len(j) != 2{
        c::__assert_fail(c"len(j) != 2", c"buildin_function.k", 8, c"main")
    }

    if typename(i) != "*i32"{
        c::__assert_fail(c"typename(i) != \"*i32\"", c"buildin_function.k", 12, c"main")
    }

    if size(i[0]) != 4{
        c::__assert_fail(c"size(i[0]) != 4", c"buildin_function.k", 16, c"main")
    }

    return 0
}