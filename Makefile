default: all

.PHONY: update-api all release release-clean

update-api:
	go get -u github.com/klev-dev/klev-api-go@main
	go mod tidy

all:
	go build -v -o klev .

release-build:
	GOOS=darwin GOARCH=amd64 go build -v -o dist/build/klev-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -v -o dist/build/klev-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build -v -o dist/build/klev-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -v -o dist/build/klev-linux-arm64 .
	GOOS=freebsd GOARCH=amd64 go build -v -o dist/build/klev-freebsd-amd64 .
	GOOS=freebsd GOARCH=arm64 go build -v -o dist/build/klev-freebsd-arm64 .
	GOOS=windows GOARCH=amd64 go build -v -o dist/build/klev-windows-amd64 .
	GOOS=windows GOARCH=arm64 go build -v -o dist/build/klev-windows-arm64 .

release: release-build
	mkdir dist/archive
	for x in $(shell ls dist/build); do tar --transform='flags=r;s|-.*||' -czf dist/archive/$$x.tar.gz -C dist/build $$x; done

release-clean:
	rm -rf dist/
