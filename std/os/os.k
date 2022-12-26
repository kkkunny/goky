import std.c

@exit
pub func exit(c: u8){
    c:__c_exit(c as i32)
}