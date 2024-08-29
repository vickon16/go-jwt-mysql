# Define the Go command
GO := go

# Define the binary output name
BINARY_NAME := go-jwt-mysql

# Define the directory for the source code
SRC_DIR := ./cmd

# Define the package for the main application
MAIN_PKG := $(SRC_DIR)/main.go

# Define the directories to format and vet
VET_DIRS := ./...

# Default target: build the application
all: build

# Build the Go binary
build:
	$(GO) build -o bin/$(BINARY_NAME) $(MAIN_PKG)

# Run Go tests
test:
	$(GO) test $(VET_DIRS)

# Run Go vet to analyze code
vet:
	$(GO) vet $(VET_DIRS)

# Format Go code
fmt:
	$(GO) fmt $(VET_DIRS)

# Clean the built binary
clean:
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GO) mod tidy

# Run the application
run: build
	./bin/$(BINARY_NAME)

# Run migrations
# e.g make migration add-user-table
migration:
	migrate create -ext sql -dir $(SRC_DIR)/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

# Run migrate up
migrate-up:
	$(GO) run $(SRC_DIR)/migrate/migrate.go up

# Run migrate down
migrate-down:
	$(GO) run $(SRC_DIR)/migrate/migrate.go down



