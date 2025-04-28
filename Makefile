BUILD_DIR := dist
BIN_NAME := aws_ipadd
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64
RELEASE_VERSION ?= dev

clean:
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)

build: clean
	@for platform in $(PLATFORMS); do \
		OS=$$(echo $$platform | cut -d'/' -f1); \
		ARCH=$$(echo $$platform | cut -d'/' -f2); \
		OUTPUT_FILE=$(BIN_NAME); \
		ARTF_FILE=$(BIN_NAME)_$${OS}_$${ARCH}; \
		echo "Building $$ARTF_FILE ..."; \
		GOOS=$$OS GOARCH=$$ARCH go build -ldflags "-X aws_ipadd/cliargs.Version=$(RELEASE_VERSION)" -o $(BUILD_DIR)/$$OUTPUT_FILE .; \
		tar -czf $(BUILD_DIR)/$$ARTF_FILE.tar.gz -C $(BUILD_DIR) $$OUTPUT_FILE; \
		rm -f $(BUILD_DIR)/$$OUTPUT_FILE; \
	done
	@echo "Build complete. Artifacts are in $(BUILD_DIR)/"
