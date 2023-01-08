pub type tm struct{}
pub type clock_t long
pub type time_t long

@extern(asctime)
pub func asctime(timeptr: *tm)*char

@extern(clock)
pub func clock()clock_t

@extern(ctime)
pub func ctime(timer: *time_t)*char

@extern(difftime)
pub func difftime(time1: time_t, time2: time_t)double

@extern(gmtime)
pub func gmtime(timer: *time_t)*tm

@extern(localtime)
pub func localtime(timer: *time_t)*tm

@extern(mktime)
pub func mktime(timeptr: *tm)time_t

@extern(strftime)
pub func strftime(str: *char, maxsize: size_t, format: *char, timeptr: *tm)size_t

@extern(time)
pub func time(timer: *time_t)time_t