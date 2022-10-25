BINARY ?= r8-retrace
VERSION ?= v1.0.1
DEFAULT_PORT ?= 8082

CMD=go build -ldflags "-X main.version=$(VERSION) -X main.defalutPort=$(DEFAULT_PORT)"
GO_BUILD=$(CMD) -o bin/$(BINARY)-$(VERSION)-$@

.PHONY:clean
clean:
	@rm -rf ./bin

.PHONY:build
build:
	$(CMD) -o bin/$(BINARY)-$(VERSION)

.PHONY:mac-x86_64
mac-x86_64:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=darwin $(GO_BUILD)

.PHONY:mac-arm64
mac-arm64:
	GOARCH=arm64 CGO_ENABLED=0 GOOS=darwin $(GO_BUILD)


.PHONY:linux-x86_64
linux-x86_64:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=linux  $(GO_BUILD)

.PHONY:linux-arm64
linux-arm64:
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux  $(GO_BUILD)


.PHONY: win-x86_64
win-x86_64.exe:
	GOARCH=amd64 CGO_ENABLED=0 GOOS=windows $(GO_BUILD)

.PHONY: win-386
win-386.exe:
	GOARCH=386 CGO_ENABLED=0  GOOS=windows $(GO_BUILD)


.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: run
run:
	@go run .

.PHONY: all
all: clean mac-x86_64 linux-x86_64 win-x86_64.exe mac-arm64 linux-arm64 win-386.exe

install-my-mac-tool: build
	mv bin/$(BINARY)-$(VERSION) ~/tools/path/$(BINARY)
	rm -r bin

help:
	@echo "make run - 直接运行 Go 代码"
	@echo "make build - 编译 Go 代码, 并编译生成二进制文件"
	@echo "make clean - 移除编译的二进制文件"
	@echo "make all - 编译多平台的二进制文件"
	@echo "make mac-amd64 - 编译 Go 代码, 生成mac-amd64的二进制文件"
	@echo "make linux-amd64 - 编译 Go 代码, 生成linux-amd64二进制文件"
	@echo "make win-x86_64 - 编译 Go 代码, 生成windows-amd64二进制文件"
	@echo "make tidy - 执行go mod tidy"
