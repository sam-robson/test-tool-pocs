COMMANDER_DIR = ./commander
INK_DIR = ./ink
BUBBLE_TEA_DIR = ./bubble-tea
BIN_DIR = ./bin

build_commander:
	cd $(COMMANDER_DIR) && npm run build
	chmod +x $(COMMANDER_DIR)/dist/index.js

build_ink:
	cd $(INK_DIR) && npm run build
	chmod +x $(INK_DIR)/dist/cli.js

build_bubble_tea:
	mkdir -p $(BIN_DIR)
	cd $(BUBBLE_TEA_DIR) && go build -o ../$(BIN_DIR)/bubble-tea test.go

all: build_commander build_ink build_bubble_tea

.PHONY: all build_commander build_ink build_bubble_tea