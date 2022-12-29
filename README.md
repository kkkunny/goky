# klang
The K programming language, a strongly typed, compiled programming language with a goal of simplicity.

##### Hello World
```go
import std.c

@extern(main)
func main()u8{
    let s = c"hello world"
    c:puts(&s as *c:char)
    return 0
}
```