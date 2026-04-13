# 项目基本信息
BINARY_NAME=kamailio-manager
MAIN_PATH=./cmd/server/main.go
DIST_DIR=./dist
CONF_SRC=./configs/config.yaml
CONF_DEST=/usr/local/kamailio/etc/kamailio/managerConf.yaml

# 编译参数
GO=go
GOFLAGS=-ldflags="-s -w" # 压缩体积，移除调试信息

.PHONY: all build clean run help setup-dir

## all: 默认目标，执行编译
all: setup-dir build

## setup-dir: 创建输出目录
setup-dir:
	@mkdir -p $(DIST_DIR)

## build: 编译当前平台的二进制文件
build:
	@echo "Building for current platform..."
	$(GO) build $(GOFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(DIST_DIR)/$(BINARY_NAME)"

## linux: 交叉编译 Linux 64位版本 (适合部署到服务器)
linux: setup-dir
	@echo "Cross-compiling for Linux AMD64..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)_linux $(MAIN_PATH)
	@echo "Cross-compile complete: $(DIST_DIR)/$(BINARY_NAME)_linux"

## run: 直接运行项目
run:
	$(GO) run $(MAIN_PATH)

## clean: 清理编译产物
clean:
	@echo "Cleaning up..."
	@rm -rf $(DIST_DIR)
	@echo "Clean complete."

## install: 将配置文件拷贝到系统目录 (需要 sudo)
install-conf:
	@echo "Installing config to $(CONF_DEST)..."
	@mkdir -p /usr/local/kamailio/etc/kamailio/
	@cp $(CONF_SRC) $(CONF_DEST)
	@echo "Config installed."

## help: 显示帮助信息
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/##//'
