BINARY_NAME=scry
INSTALL_DIR=$(HOME)/.local/bin
ROOT_DIR=$(HOME)/.local/share/scry
DATA_DIR=$(ROOT_DIR)/data
INDEX_DIR=$(ROOT_DIR)/index

.PHONY: all build install clean uninstall purge

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o bin/$(BINARY_NAME) ./cmd/scry

install: build
	@echo "Installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@mkdir -p $(DATA_DIR)
	@mkdir -p $(INDEX_DIR)
	@cp bin/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Scry installed successfully!"

uninstall:
	@echo "Removing $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	@echo "Cleaning build directory..."
	@rm -rf bin/

purge: uninstall
	@echo "Removing all downloaded indexes and data..."
	@rm -rf $(ROOT_DIR)
	@echo "Scry data completely wiped."
