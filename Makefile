.PHONY: all build build-load build-pv build-wind docker-build docker-build-load docker-build-pv docker-build-wind clean help

# Default target
all: build docker-build

# Build all Go binaries
build: build-load build-pv build-wind

# Build individual Go binaries
build-load:
	@echo "Building asset-sim-load..."
	cd asset-sim-load && go build -o bin/asset-sim-load main.go

build-pv:
	@echo "Building asset-sim-pv..."
	cd asset-sim-pv && go build -o bin/asset-sim-pv main.go

build-wind:
	@echo "Building asset-sim-wind..."
	cd asset-sim-wind && go build -o bin/asset-sim-wind main.go

# Build all Docker images
docker-build: docker-build-load docker-build-pv docker-build-wind

# Build individual Docker images
docker-build-load:
	@echo "Building Docker image for asset-sim-load..."
	cd asset-sim-load && $(MAKE) docker-build

docker-build-pv:
	@echo "Building Docker image for asset-sim-pv..."
	cd asset-sim-pv && $(MAKE) docker-build

docker-build-wind:
	@echo "Building Docker image for asset-sim-wind..."
	cd asset-sim-wind && $(MAKE) docker-build

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf asset-sim-load/bin asset-sim-load/asset-sim-load
	rm -rf asset-sim-pv/bin asset-sim-pv/asset-sim-pv
	rm -rf asset-sim-wind/bin asset-sim-wind/asset-sim-wind

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Build all Go binaries and Docker images (default)"
	@echo "  build            - Build all Go binaries"
	@echo "  build-load       - Build asset-sim-load Go binary"
	@echo "  build-pv         - Build asset-sim-pv Go binary"
	@echo "  build-wind       - Build asset-sim-wind Go binary"
	@echo "  docker-build     - Build all Docker images"
	@echo "  docker-build-load - Build asset-sim-load Docker image"
	@echo "  docker-build-pv   - Build asset-sim-pv Docker image"
	@echo "  docker-build-wind - Build asset-sim-wind Docker image"
	@echo "  clean            - Remove all build artifacts"
	@echo "  help             - Show this help message"
