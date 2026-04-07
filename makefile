all: run
build:
	sh create_disk.sh
	@go build -o bin/kdfs
run: build
	@./bin/kdfs
test:
	@go test ./ -v