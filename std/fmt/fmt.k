import std.sys

pub func print(s: *i8){
    let length: usize = 0
    for true{
        let c = s[length]
        if c == 0{
            break
        }
        length += 1
    }
    sys:write(1, s, length)
}