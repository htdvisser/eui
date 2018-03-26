.PHONY: dev-deps

dev-deps:
	go get -u golang.org/x/vgo

.PHONY: build

build:
	vgo build -o dist/whois-eui-$(shell go env GOOS)-$(shell go env GOARCH)$(shell go env GOEXE) cmd/whois-eui/main.go
