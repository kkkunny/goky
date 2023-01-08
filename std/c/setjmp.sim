pub type __jmp_buf_tag struct{}
pub type jmp_buf [1]__jmp_buf_tag

@extern(setjmp)
pub func setjmp(env: jmp_buf)int

@extern(longjmp)
@noreturn
pub func longjmp(env: jmp_buf, val: int)