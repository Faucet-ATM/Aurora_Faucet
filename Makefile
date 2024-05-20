.PHONY: build run clean clean-cache install

install:
	@echo "Installing dependencies..."
	@go mod tidy

build:
	@echo "Building the project..."
	@go build -o build/aurora-faucet ./cmd/aurora-faucet

run: build
	@echo "Running the project..."
	@./build/aurora-faucet

clean:
	@echo "Cleaning up build files..."
	@rm -rf build

clean-cache:
	@echo "Cleaning up build cache..."
	@go clean -cache -modcache -i -r
