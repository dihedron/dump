.PHONY: build
build:
	@ENABLE_CGO=0 go build -o bin/microservice
	@cp configuration.yaml bin/

.PHONY: clean
clean:
	@rm -rf bin/

.PHONY: run-client-on-9080
run-server-on-9080: build
	bin/microservice --port=9080 create account1 1000