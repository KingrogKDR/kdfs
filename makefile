all: run

build-disk:
	@sh create_disk.sh
build:
	@go build -o bin/kdfs
run: build
	@./bin/kdfs
test:
	@go test ./... -v