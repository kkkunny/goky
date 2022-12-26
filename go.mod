module github.com/kkkunny/klang

go 1.19

require (
	github.com/alecthomas/participle/v2 v2.0.0-beta.5
	github.com/kkkunny/stl v0.0.0-20221015140421-2a6594a9d191
	github.com/spf13/cobra v1.6.1
	golang.org/x/exp v0.0.0-20220414153411-bcd21879b8fd
	tinygo.org/x/go-llvm v0.0.0-20221212185523-e80bc424a2b1
)

require (
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

replace tinygo.org/x/go-llvm v0.0.0-20221212185523-e80bc424a2b1 => /mnt/code/go/src/github.com/kkkunny/go-llvm
