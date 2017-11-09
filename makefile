.PHONY: build clean test run default

build: clean
	@go build -o ./revisor ./cmd/revisor

clean:
	@rm -rf ./revisor

test:
	go test ./...

run: build
	@./revisor -root ./public/ -base https://www.example.com/

default: build
