# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=gin-basic
BINARY_UNIX=gin-basic-unix

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v

# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Test the project
test:
	$(GOTEST) -v ./...

# Build for unix systems
build-linux:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BINARY_UNIX) -v

# Docker build
docker-build:
	docker build -t gin-basic .

# Docker run
docker-run:
	docker run -p 8080:8080 gin-basic

# Install dependencies
deps:
	go mod tidy

# Start the application in background
start:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME) &

# Stop the application
stop:
	pkill -f $(BINARY_NAME) || true

.PHONY: build run clean test build-linux docker-build docker-run deps start stop