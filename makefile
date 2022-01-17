test_unit:
	go test -count=1 -race -cover -short  ./...

test:
	go test -count=1 -race -cover  ./...

build:
	go build -o server

run:
	export config_file=config.yml && ./server serve

test_run:
	export config_file=config.yml && go run main.go serve
