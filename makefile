# Go parameters
GOCMD=go
PORT?=8080  # Default port, can be overridden
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=blog-server
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server

# Templ parameters
TEMPL=templ
TEMPL_DIR=./web
TEMPL_FILES=$(shell find $(TEMPL_DIR) -name "*.templ")

# Tailwind parameters
NPX=npx
TAILWIND=tailwindcss
CSS_INPUT=./web/styles/input.css
CSS_OUTPUT=./web/static/css/main.css

# Air for hot reloading
AIR=air

.PHONY: all build clean test coverage deps air run generate-css watch-css generate-templ watch-templ dev install-tools help

all: clean deps build

# Build the application
build: generate-templ generate-css
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf ./tmp

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Download all dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	npm install

# Run with air (hot reloading)
air:
	$(AIR)

# Run the application
run: build
	PORT=$(PORT) ./$(BINARY_NAME)

# Run server on specific port
serve:
	PORT=$(PORT) $(GORUN) $(MAIN_PATH)/main.go

# Generate CSS
generate-css:
	$(NPX) $(TAILWIND) -i $(CSS_INPUT) -o $(CSS_OUTPUT) --minify

# Watch CSS changes
watch-css:
	$(NPX) $(TAILWIND) -i $(CSS_INPUT) -o $(CSS_OUTPUT) --watch

# Generate Templ files
generate-templ:
	$(TEMPL) generate

# Watch Templ files
watch-templ:
	$(TEMPL) generate --watch

# Development mode (run all watchers)
dev:
	make generate-templ
	make generate-css
	$(AIR)

# Install required tools
install-tools:
	$(GOGET) -u github.com/a-h/templ/cmd/templ
	$(GOGET) -u github.com/cosmtrek/air
	npm install -D tailwindcss
	npm install -D @tailwindcss/typography
	npm install -D @tailwindcss/forms
	npm install alpinejs
	npm install htmx.org

# Create a new blog post template
new-post:
	@read -p "Enter post title: " title; \
	filename=`echo $$title | tr '[:upper:]' '[:lower:]' | tr ' ' '-'`; \
	date=`date +%Y-%m-%d`; \
	echo "Creating new post: $$filename"; \
	echo "---\ntitle: $$title\ndate: $$date\ndraft: true\n---\n\n" > ./content/posts/$$filename.md

# Build for production
build-prod: generate-templ
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Database commands
.PHONY: db-setup db-migrate db-rollback

# Setup database
db-setup:
	mkdir -p data migrations

# Run migrations
db-migrate:
	@echo "Running database migrations..."
	@go run cmd/server/main.go migrate

# Rollback last migration
db-rollback:
	@echo "Rolling back last migration..."
	@go run cmd/server/main.go rollback
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build files"
	@echo "  make test           - Run tests"
	@echo "  make coverage       - Run tests with coverage"
	@echo "  make deps           - Download Go dependencies"
	@echo "  make air            - Run with hot reloading"
	@echo "  make run            - Run the application"
	@echo "  make generate-css   - Generate CSS files"
	@echo "  make watch-css      - Watch CSS changes"
	@echo "  make generate-templ - Generate Templ files"
	@echo "  make watch-templ    - Watch Templ files"
	@echo "  make dev            - Run in development mode"
	@echo "  make install-tools  - Install required tools"
	@echo "  make new-post      - Create a new blog post"
	@echo "  make build-prod    - Build for production"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
