VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION) -s -w"
BIN      := bin/irecall
MCP_BIN  := bin/irecall-mcp
WEB_BIN  := bin/irecall-web
WEB_WINDOWS_BIN := bin/irecall-web-windows-amd64.exe
DESKTOP_BIN := bin/irecall-desktop
DESKTOP_WINDOWS_BIN := bin/irecall-desktop-windows-amd64.exe
FRONTEND_DIR := frontend
WAILS_BUILD_TAGS := wails,production

.PHONY: build build-cli build-mcp build-web build-web-windows build-desktop build-desktop-windows build-local build-everything frontend-install frontend-build test lint install clean run tidy

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

run: build
	./$(BIN)

test:
	go test ./...

lint:
	go vet ./...

tidy:
	go mod tidy

install:
	go install $(LDFLAGS) ./cmd/irecall

clean:
	rm -rf bin/ $(FRONTEND_DIR)/dist

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
