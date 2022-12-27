pub type size_t usize
pub type FILE struct{}

@extern(exit)
@exit
pub func exit(i32)

@extern(puts)
pub func puts(*i8)