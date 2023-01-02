import std.c

@extern(main)
func main()u8{
    let fileStr = c"buildin_function.k"
    let funcStr = c"main"

    let i = "你好世界！"
    let str1 = c"len(i) != 5"
    if len(i) != 5{
        c::__assert_fail(&str1 as *c::char, &fileStr as *c::char, 11, &funcStr as *c::char)
    }

    let str2 = c"typename(i) != \"[5]i32\""
    if typename(i) != "[5]i32"{
        c::__assert_fail(&str2 as *c::char, &fileStr as *c::char, 16, &funcStr as *c::char)
    }

    let str3 = c"size(i[0]) != 4"
    if size(i[0]) != 4{
        c::__assert_fail(&str3 as *c::char, &fileStr as *c::char, 21, &funcStr as *c::char)
    }

    return 0
}