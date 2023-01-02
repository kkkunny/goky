import std.c

// 退出程序
@exit
@noreturn
pub func exit(code: u8){
    c::exit(code as c::int)
}