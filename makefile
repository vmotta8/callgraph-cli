BIN_NAME = callgraph-cli
RUST_CLI_DIR = clis/rust

.PHONY: build analyze-cg test clean help build-go build-rust clean-go clean-rust tool-version

build: build-go build-rust

build-go:
	@echo "Building Go CLI..."
	@go build -o bin/$(BIN_NAME) ./cmd

build-rust:
	@echo "Building Rust CLI..."
	@cd $(RUST_CLI_DIR) && cargo build --release

clean: clean-go clean-rust

clean-go:
	@echo "Cleaning Go binaries..."
	@rm -rf bin/

clean-rust:
	@echo "Cleaning Rust target directory..."
	@cd $(RUST_CLI_DIR) && cargo clean

test:
	@go test -v ./...

analyze-cg:
	./bin/$(BIN_NAME) analyze-cg -e $(entryFile) -f $(funcName) -s $(stdout)
