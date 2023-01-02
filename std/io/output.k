import std.c
import std.container.string

pub func print(s: string::String){
    let index: usize
    for index < s.len{
        c::putwchar(s.data[index] as c::wchar_t)
        index += 1
    }
}