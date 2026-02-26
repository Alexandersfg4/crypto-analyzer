app_name = crypto-analyzer

run:
	go run cmd/main.go

build:
	go build cmd/main.go -o bin/$(app_name)

install:
	cp bin/$(app_name) /usr/local/bin/
