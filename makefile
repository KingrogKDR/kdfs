all: run
build:
	@go build -gcflags="-m" -o bin/kdfs
run: build
	@./bin/kdfs
test:
	@go test ./ -v