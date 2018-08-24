.PHONY: build

build:
	go build -o dist/whois-eui-$(shell go env GOOS)-$(shell go env GOARCH)$(shell go env GOEXE) ./cmd/whois-eui
