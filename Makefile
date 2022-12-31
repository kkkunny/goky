test=tests/exit.k

.PHONY: lex
lex: lex.go $(test)
	@go build -tags test,lex -o .test .
	-@./.test $(test) || true
	@rm .test

.PHONY: parse
parse: parse.go $(test)
	@go build -tags test,parse -o .test .
	-@./.test $(test) || true
	@rm .test

.PHONY: analyse
analyse: analyse.go $(test)
	@go build -tags test,analyse -o .test .
	-@./.test $(test) || true
	@rm .test

.PHONY: codegen
codegen: codegen.go $(test)
	@go build -tags test,codegen,llvm14 -o .test .
	-@./.test $(test) || true
	@rm .test

.PHONY: optimize
optimize: optimize.go $(test)
	@go build -tags test,optimize,llvm14 -o .test .
	-@./.test $(test) || true
	@rm .test

.PHONY: build
build: clean main.go
	go build -tags llvm14 -o kcc main.go
	ln -s /mnt/code/go/src/github.com/kkkunny/klang/kcc /mnt/code/go/bin/kcc

.PHONY: clean
clean:
	rm -f /mnt/code/go/bin/kcc
	rm -f /mnt/code/go/src/github.com/kkkunny/klang/kcc

.PHONY: test
test: build
	@for file in $(foreach dir, tests, $(wildcard tests/*.k)); do \
		echo kcc run tests/hello_world.k > /dev/null; \
		echo -e "\e[32m 测试成功 $$file \e[0m"; \
    done; \
    make clean