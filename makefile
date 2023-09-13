test_unit:
	go test -count=1 -race -cover -short  ./...

test:
	go test -count=1 -race -cover  ./...

build:
	go build -o server

run: build
	export CONFIG_FILE=config.yaml && ./server

dev:
	export CONFIG_FILE=config-dev.yaml && go run main.go serve
