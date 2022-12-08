@exit
@extern(exit)
func __c_exit(i32)

@exit
pub func exit(code: u8){
    __c_exit(code as i32)
}
