// 字符串
pub type String struct {
    data: *i32
    len: usize
}

pub func new(s: *i32) String {
    let len: usize
    for s[len] != 0 {
        len += 1
    }
    return {s, len}
}