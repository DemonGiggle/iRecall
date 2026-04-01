VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION) -s -w"
BIN      := bin/irecall

.PHONY: build test lint install clean run tidy

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o $(BIN) ./cmd/irecall

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
	rm -rf bin/

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
