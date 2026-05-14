BIN := purr
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: build install uninstall test clean

build:
	go build -o $(BIN) ./cmd/purr

install:
	@mkdir -p $(INSTALL_DIR)
	go build -o $(INSTALL_DIR)/$(BIN) ./cmd/purr
	@echo "installed → $(INSTALL_DIR)/$(BIN)"

uninstall:
	rm -f $(INSTALL_DIR)/$(BIN)

test:
	go test ./...

clean:
	rm -f $(BIN)
