    .intel_syntax noprefix
    .file	"syscall.s"
    .text
    .globl	__syscall
    .type	__syscall,@function
__syscall:
    enter 0, 0
    mov rax, rdi
    mov rdi, rsi
    mov rsi, rdx
    mov rdx, rcx
    mov r10, r8
    mov r8, r9
    mov r9, qword ptr [rbp + 16]
    syscall
    leave
    ret

