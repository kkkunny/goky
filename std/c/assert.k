@exit
@noreturn
@extern(__assert_fail)
pub func __assert_fail(assertion: *char, file: *char, line: unsigned_int, function: *char)

@exit
@noreturn
@extern(__assert_perror_fail)
pub func __assert_perror_fail(errnum: int, file: *char, line: unsigned_int, function: *char)

@exit
@noreturn
@extern(__assert)
pub func __assert(assertion: *char, file: *char, line: int)