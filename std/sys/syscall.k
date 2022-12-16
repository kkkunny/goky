@extern(__syscall)
@link(asm="syscall.s")
pub func syscall(trap: usize, a1: usize, a2: usize, a3: usize, a4: usize, a5: usize, a6: usize)

pub func write(fd: usize, buf: *i8, count: usize){
    syscall(1, fd, buf as usize, count, 0, 0, 0)
}

@exit
pub func exit(code: u8){
    syscall(60, code as usize, 0, 0, 0, 0, 0)
}