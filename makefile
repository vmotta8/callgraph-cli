BIN_NAME = callgraph-cli
RUST_CLI_DIR = clis/rust
RUST_BIN_NAME = rust-callgraph-cli
RUST_DEST_DIR = internal/languages/rust

.PHONY: install build build-go build-rust clean clean-go clean-rust test analyze-cg

install: build
	@echo "Installing $(BIN_NAME) into /usr/local/bin..."
	@sudo cp bin/$(BIN_NAME) /usr/local/bin/

build: build-rust build-go

build-go:
	@echo "Building Go CLI..."
	@go build -o bin/$(BIN_NAME) ./cmd

build-rust:
	@echo "Building Rust CLI..."
	@cd $(RUST_CLI_DIR) && cargo build --release
	@echo "Copying Rust binary to $(RUST_DEST_DIR)..."
	@mkdir -p $(RUST_DEST_DIR)
	@cp $(RUST_CLI_DIR)/target/release/$(RUST_BIN_NAME) $(RUST_DEST_DIR)/

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
