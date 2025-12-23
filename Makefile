.PHONY: lint format test build clean

# Run all linters
lint:
	@echo "Running go fmt..."
	go fmt ./...
	@echo "Running go vet..."
	go vet ./...
	@echo "Running goimports..."
	@which goimports > /dev/null || (echo "Installing goimports..." && go install golang.org/x/tools/cmd/goimports@latest)
	goimports -w .
	@echo "Running staticcheck..."
	@which staticcheck > /dev/null || (echo "Installing staticcheck..." && go install honnef.co/go/tools/cmd/staticcheck@latest)
	staticcheck ./...
	@echo "Running errcheck..."
	@which errcheck > /dev/null || (echo "Installing errcheck..." && go install github.com/kisielk/errcheck@latest)
	errcheck ./...

# Format code
format:
	go fmt ./...
	goimports -w .

# Run tests
test:
	go test -v ./...

# Build the application
build:
	go build -o bin/server ./cmd/api

# Run the application
run:
	go run ./cmd/api/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Install all linting tools
install-lint-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/kisielk/errcheck@latest

# Run golangci-lint with working config
ci-lint:
	$(GOPATH)/bin/golangci-lint run --disable-all --enable=gofmt,goimports,govet,errcheck,staticcheck,unused
