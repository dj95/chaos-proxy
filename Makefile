# Go settings
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOLINT=golint
GOGET=$(GOCMD) get

# Build settings
BINARY_PATH=./bin/
BINARY_NAME=deception-proxy
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_WINDOWS=$(BINARY_NAME).exe
BINARY_MACOS=$(BINARY_NAME)_macos

# Test Settings
TEST_FILES := $(shell $(GOCMD) list ./...)

all: deps test

run:
		$(GOBUILD) -o $(BINARY_NAME) -v cmd/deception-proxy/main.go
		./$(BINARY_NAME) --config configs/config.yml

build:
		$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) -v cmd/deception-proxy/main.go

test:
		mkdir -p report
		$(GOTEST) -v -short -covermode=count -coverprofile report/cover.out $(TEST_FILES)
		$(GOCMD) tool cover -html=report/cover.out -o report/cover.html
		$(GOLINT) -set_exit_status $(TEST_FILES)
		CC=clang $(GOTEST) -v -msan -short $(TEST_FILES)
		staticcheck $(TEST_FILES)

clean:
		$(GOCLEAN)
		rm -f $(BINARY_PATH)$(BINARY_NAME)

deps:
		GO111MODULE=on $(GOCMD) mod vendor
