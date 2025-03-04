BUILD_DIR := dist
BIN_NAME := aws_ipadd
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64

clean:
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)

build: clean
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_FILE=$(BIN_NAME)_$${OS}_$${ARCH}; \
		echo "Building $$OUTPUT_FILE..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $(BUILD_DIR)/$$OUTPUT_FILE .; \
		tar -czf $(BUILD_DIR)/$$OUTPUT_FILE.tar.gz -C $(BUILD_DIR) $$OUTPUT_FILE; \
		rm -f $(BUILD_DIR)/$$OUTPUT_FILE; \
	done
	@echo "Build complete. Artifacts are in $(BUILD_DIR)/"
