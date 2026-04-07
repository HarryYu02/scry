BINARY_NAME=scry
INSTALL_DIR=$(HOME)/.local/bin
SCRY_DIR=$(HOME)/.local/share/scry
SCRY_DATA_DIR=$(SCRY_DIR)/data
SCRY_INDEX_DIR=$(SCRY_DIR)/index

.PHONY: all build install clean uninstall purge

all: build

build:
	@echo "building $(BINARY_NAME)..."
	@go build -ldflags="-s -w" -o bin/$(BINARY_NAME) ./cmd/scry
	@echo "scry built successfully!"

install: build
	@echo "installing to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@mkdir -p $(SCRY_DATA_DIR)
	@mkdir -p $(SCRY_INDEX_DIR)
	@cp bin/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "scry installed successfully!"

uninstall:
	@echo "removing $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "scry removed!"

clean:
	@echo "cleaning build directory..."
	@rm -rf bin/
	@echo "all build directory cleaned!"

purge: uninstall
	@echo "removing all downloaded indexes and data..."
	@rm -rf $(SCRY_DIR)
	@echo "scry data completely wiped."
