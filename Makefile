GO_BINARY_NAME=go-test
BIN_DIR=$(PWD)/bin

.PHONY: build clean

goinit:
	go mod init test-tool
	go get github.com/charmbracelet/bubbletea
	go get github.com/charmbracelet/lipgloss

gobuild:
	go build -o $(BIN_DIR)/$(GO_BINARY_NAME) bubble-tea/test.go

gorun: gobuild
	$(BIN_DIR)/$(GO_BINARY_NAME)

install: gobuild
	mv $(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BIN_DIR)