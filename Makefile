.PHONY: build
build:
	@ENABLE_CGO=0 go build -o bin/dump
	@cp configuration.yaml bin/

.PHONY: clean
clean:
	@rm -rf bin/

.PHONY: run
run: build
	bin/dump --port=9080 create account1 1000