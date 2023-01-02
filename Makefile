WORK_PATH = $(shell pwd)
TEST_DIR = $(WORK_PATH)/tests
TEST_FILE = $(TEST_DIR)/hello_world.k
KCC_FILE = $(GOPATH)/bin/kcc

.PHONY: lex
lex: lex.go $(TEST_FILE)
	@go build -tags test,lex -o .test .
	-@./.test $(TEST_FILE) || true
	@rm .test

.PHONY: parse
parse: parse.go $(TEST_FILE)
	@go build -tags test,parse -o .test .
	-@./.test $(TEST_FILE) || true
	@rm .test

.PHONY: analyse
analyse: analyse.go $(TEST_FILE)
	@go build -tags test,analyse -o .test .
	-@./.test $(TEST_FILE) || true
	@rm .test

.PHONY: codegen
codegen: codegen.go $(TEST_FILE)
	@go build -tags test,codegen,llvm14 -o .test .
	-@./.test $(TEST_FILE) || true
	@rm .test

.PHONY: optimize
optimize: optimize.go $(TEST_FILE)
	@go build -tags test,optimize,llvm14 -o .test .
	-@./.test $(TEST_FILE) || true
	@rm .test

.PHONY: build
build: clean main.go
	go build -tags llvm14 -o kcc main.go
	ln -s $(WORK_PATH)/kcc $(KCC_FILE)

.PHONY: clean
clean:
	rm -f $(WORK_PATH)/kcc
	rm -f $(KCC_FILE)

.PHONY: test
test: build
	@for file in $(foreach dir, $(TEST_DIR), $(wildcard $(TEST_DIR)/*.k)); do \
		echo kcc run $$file > /dev/null; \
		echo -e "\e[32m 测试成功 $$file \e[0m"; \
    done; \
    make clean