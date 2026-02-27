app_name = crypto-analyzer

run:
	go run cmd/main.go

build:
	go build -o bin/$(app_name) cmd/main.go

install:
	cp bin/$(app_name) /usr/local/bin/
