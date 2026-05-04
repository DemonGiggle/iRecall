VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION) -s -w"
BIN      := bin/irecall
MCP_BIN  := bin/irecall-mcp
WEB_BIN  := bin/irecall-web
WEB_WINDOWS_BIN := bin/irecall-web-windows-amd64.exe
DESKTOP_BIN := bin/irecall-desktop
DESKTOP_WINDOWS_BIN := bin/irecall-desktop-windows-amd64.exe
FRONTEND_DIR := frontend
WAILS_BUILD_TAGS := wails,production,$(shell pkg-config --exists webkit2gtk-4.1 2>/dev/null && echo ,webkit2_41)
RELEASE_DIR := dist/release
RELEASE_TUI_TARGETS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64
RELEASE_WEB_SERVER_TARGETS := $(RELEASE_TUI_TARGETS)
# Wails desktop cross-builds are currently limited to the targets validated by this repo's toolchain.
RELEASE_DESKTOP_TARGETS ?= linux/amd64 windows/amd64 windows/arm64

.PHONY: build build-cli build-mcp build-web build-web-windows build-desktop build-desktop-windows build-local build-everything build-release clean-release frontend-install frontend-build test test-mcp-bootstrap lint install clean run tidy

build: build-cli

build-cli:
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BIN) ./cmd/irecall

build-mcp:
	@mkdir -p bin
	go build $(LDFLAGS) -o $(MCP_BIN) ./cmd/irecall-mcp

build-web: frontend-build
	@mkdir -p bin
	go build $(LDFLAGS) -o $(WEB_BIN) ./web

build-web-windows: frontend-build
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(WEB_WINDOWS_BIN) ./web

frontend-install:
	cd $(FRONTEND_DIR) && if [ -f package-lock.json ]; then rm -rf node_modules && npm ci; else npm install; fi

frontend-build: frontend-install
	cd $(FRONTEND_DIR) && npm run build

build-desktop: frontend-build
	@mkdir -p bin
	go build -tags "$(WAILS_BUILD_TAGS)" -o $(DESKTOP_BIN) ./desktop

build-desktop-windows: frontend-build
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build -tags "$(WAILS_BUILD_TAGS)" -o $(DESKTOP_WINDOWS_BIN) ./desktop

build-local: build-cli build-web build-desktop

build-everything: build-local build-all

build-release: clean-release frontend-build
	@mkdir -p $(RELEASE_DIR)
	@set -eu; \
	for target in $(RELEASE_TUI_TARGETS); do \
		os=$${target%/*}; \
		arch=$${target#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		stage="$(RELEASE_DIR)/irecall-tui-$(VERSION)-$$os-$$arch"; \
		mkdir -p "$$stage"; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o "$$stage/irecall$$ext" ./cmd/irecall; \
		tar -C "$(RELEASE_DIR)" -czf "$(RELEASE_DIR)/irecall-tui-$(VERSION)-$$os-$$arch.tar.gz" "irecall-tui-$(VERSION)-$$os-$$arch"; \
		rm -rf "$$stage"; \
	done; \
	for target in $(RELEASE_WEB_SERVER_TARGETS); do \
		os=$${target%/*}; \
		arch=$${target#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		stage="$(RELEASE_DIR)/irecall-web-server-$(VERSION)-$$os-$$arch"; \
		mkdir -p "$$stage"; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o "$$stage/irecall-web$$ext" ./web; \
		tar -C "$(RELEASE_DIR)" -czf "$(RELEASE_DIR)/irecall-web-server-$(VERSION)-$$os-$$arch.tar.gz" "irecall-web-server-$(VERSION)-$$os-$$arch"; \
		rm -rf "$$stage"; \
	done; \
	for target in $(RELEASE_DESKTOP_TARGETS); do \
		os=$${target%/*}; \
		arch=$${target#*/}; \
		ext=""; \
		if [ "$$os" = "windows" ]; then ext=".exe"; fi; \
		stage="$(RELEASE_DIR)/irecall-desktop-$(VERSION)-$$os-$$arch"; \
		mkdir -p "$$stage"; \
		GOOS=$$os GOARCH=$$arch go build -tags "$(WAILS_BUILD_TAGS)" -o "$$stage/irecall-desktop$$ext" ./desktop; \
		tar -C "$(RELEASE_DIR)" -czf "$(RELEASE_DIR)/irecall-desktop-$(VERSION)-$$os-$$arch.tar.gz" "irecall-desktop-$(VERSION)-$$os-$$arch"; \
		rm -rf "$$stage"; \
	done; \
	stage="$(RELEASE_DIR)/irecall-web-$(VERSION)"; \
	mkdir -p "$$stage"; \
	cp -R "$(FRONTEND_DIR)/dist" "$$stage/"; \
	tar -C "$(RELEASE_DIR)" -czf "$(RELEASE_DIR)/irecall-web-$(VERSION).tar.gz" "irecall-web-$(VERSION)"; \
	rm -rf "$$stage"; \
	checksum_cmd="shasum -a 256"; \
	if command -v sha256sum >/dev/null 2>&1; then checksum_cmd="sha256sum"; fi; \
	cd "$(RELEASE_DIR)" && $$checksum_cmd *.tar.gz > SHA256SUMS

run: build
	./$(BIN)

test:
	go test ./...

test-mcp-bootstrap:
	go test ./web -run TestOperatorBootstrapIssuesTokenAndMCPHealthChecksRealWebServer -count=1

lint:
	go vet ./...

tidy:
	go mod tidy

install:
	go install $(LDFLAGS) ./cmd/irecall

clean:
	rm -rf bin/ $(FRONTEND_DIR)/dist $(RELEASE_DIR)

clean-release:
	rm -rf $(RELEASE_DIR)

# Cross-compilation targets
build-linux-amd64:
	GOOS=linux   GOARCH=amd64  go build $(LDFLAGS) -o bin/irecall-linux-amd64  ./cmd/irecall

build-linux-arm64:
	GOOS=linux   GOARCH=arm64  go build $(LDFLAGS) -o bin/irecall-linux-arm64  ./cmd/irecall

build-darwin-amd64:
	GOOS=darwin  GOARCH=amd64  go build $(LDFLAGS) -o bin/irecall-darwin-amd64 ./cmd/irecall

build-darwin-arm64:
	GOOS=darwin  GOARCH=arm64  go build $(LDFLAGS) -o bin/irecall-darwin-arm64 ./cmd/irecall

build-windows-amd64:
	GOOS=windows GOARCH=amd64  go build $(LDFLAGS) -o bin/irecall-windows-amd64.exe ./cmd/irecall

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
