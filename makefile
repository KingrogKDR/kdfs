all: run
build:
	@go build -o bin/kdfs
run: build
	@./bin/kdfs
test:
	@go test ./ -v