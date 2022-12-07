.PHONY: update-api all

update-api:
	go get -u github.com/klev-dev/klev-api-go@main
	go mod tidy

all:
	go build -v -o klev .
