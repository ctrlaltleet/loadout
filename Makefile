PROJECT_NAME := loadout
BUILD_DIR := build
GOFLAGS := -ldflags "-s -w" -trimpath -buildvcs=false
GO_BUILD := go build $(GOFLAGS)
CMD_PATH := ./cmd/loadout/
.PHONY: all clean linux windows darwin tidy

all: tidy linux windows darwin

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

tidy:
	go mod tidy

linux: linux-amd64 linux-arm64

linux-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64

$(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64: tidy | $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-amd64 $(CMD_PATH)

linux-arm64: $(BUILD_DIR)/$(PROJECT_NAME)-linux-arm64

$(BUILD_DIR)/$(PROJECT_NAME)-linux-arm64: tidy | $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-arm64 $(CMD_PATH)

windows: windows-amd64 windows-arm64

windows-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe

$(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe: tidy | $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-amd64.exe $(CMD_PATH)

windows-arm64: $(BUILD_DIR)/$(PROJECT_NAME)-windows-arm64.exe

$(BUILD_DIR)/$(PROJECT_NAME)-windows-arm64.exe: tidy | $(BUILD_DIR)
	GOOS=windows GOARCH=arm64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-arm64.exe $(CMD_PATH)

darwin: darwin-amd64 darwin-arm64

darwin-amd64: $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64

$(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64: tidy | $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-amd64 $(CMD_PATH)

darwin-arm64: $(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64

$(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64: tidy | $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-arm64 $(CMD_PATH)

clean:
	rm -rf $(BUILD_DIR)
