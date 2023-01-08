NAME=hue-picker
VERSION=`git describe --tag`
FLAGS=-ldflags="-s -w -X 'main.Version=${VERSION}'" -trimpath

.PHONY: build clean

build: build-setup build-linux build-darwin

build-setup: clean
	@mkdir -p dist
	@go generate ./...

build-linux: build-linux-amd64 build-linux-arm64 build-linux-armv7

build-linux-amd64:
	@echo "Building Linux-x86_64"
	@GOOS=linux GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Linux-x86_64" .
	@md5 -r "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-x86_64" >> "dist/${NAME}-sha256sum.txt"

build-linux-arm64:
	@echo "Building Linux-aarch64"
	@GOOS=linux GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Linux-aarch64" .
	@md5 -r "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-aarch64" >> "dist/${NAME}-sha256sum.txt"

build-linux-armv7:
	@echo "Building Linux-armv7l"
	@GOOS=linux GOARCH=arm GOARM=7 go build ${FLAGS} -o "dist/${NAME}-Linux-armv7l" .
	@md5 -r "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Linux-armv7l" >> "dist/${NAME}-sha256sum.txt"

build-darwin: build-darwin-amd64 build-darwin-arm64

build-darwin-amd64:
	@echo "Building Darwin-x86_64"
	@GOOS=darwin GOARCH=amd64 go build ${FLAGS} -o "dist/${NAME}-Darwin-x86_64" .
	@md5 -r "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-x86_64" >> "dist/${NAME}-sha256sum.txt"

build-darwin-arm64:
	@echo "Building Darwin-arm64"
	@GOOS=darwin GOARCH=arm64 go build ${FLAGS} -o "dist/${NAME}-Darwin-arm64" .
	@md5 -r "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-md5sum.txt"
	@shasum -a 256 "dist/${NAME}-Darwin-arm64" >> "dist/${NAME}-sha256sum.txt"

clean:
	@echo "Cleaning dist"
	@rm -rf dist
