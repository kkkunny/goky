pub type __sighandler_t func(int)

@extern(signal)
pub func signal(sig: int, handler: __sighandler_t)__sighandler_t

@extern(raise)
pub func raise(sig: int)int