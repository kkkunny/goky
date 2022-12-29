.PHONY: parse
parse: parse.go main.k
	@go build -tags test,parse -o .test .
	-@./.test || true
	@rm .test

.PHONY: analyse
analyse: analyse.go main.k
	@go build -tags test,analyse -o .test .
	-@./.test || true
	@rm .test

.PHONY: codegen
codegen: codegen.go main.k
	@go build -tags test,codegen -o .test .
	-@./.test || true
	@rm .test

.PHONY: optimize
optimize: optimize.go main.k
	@go build -tags test,optimize -o .test .
	-@./.test || true
	@rm .test

.PHONY: build
build: clean main.go
	go build -o kcc main.go
	ln -s /mnt/code/go/src/github.com/kkkunny/klang/kcc /mnt/code/go/bin/kcc

.PHONY: clean
clean:
	rm -f /mnt/code/go/bin/kcc
	rm -f /mnt/code/go/src/github.com/kkkunny/klang/kcc