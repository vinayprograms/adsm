SHELL := /bin/bash 

# 'make' without arguments will show help
.DEFAULTGOAL := help

.PHONY: help all build tests report clean

BINARY:=adsm
BIN_DIR:=./bin
SRC_DIR:=./src
ROOT_DIR:=$(shell dirname $(MAKEFILE_LIST) | xargs)
TEST_PATH := $(shell sed -e 's/ /\\\ /g' <<< "$(ROOT_DIR)/test")
TEST_PACKAGES := "args,securitymodel/loaders,securitymodel/yamlmodel,securitymodel/objmodel,securitymodel/diagram"

help: # Show this help
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

all: tests build clean

build: # Build program
	@printf "Checking if %s directory exists..." $(BIN_DIR)
	@if [ ! -d $(BIN_DIR) ]; then mkdir $(BIN_DIR); printf "CREATED\n"; else printf "YES\n"; fi

	@printf "Building '%s/%s'...\n" $(BIN_DIR) $(BINARY)
	@cd src/main && go get args
	@cd $(SRC_DIR)/main && go build -ldflags "-s -w" -v -o ../../$(BIN_DIR)/$(BINARY) .
	@printf "DONE\n"

tests: # Run automated tests
	@printf "Running automated tests...\n"

	@cd $(TEST_PATH) && go get args && go get -t tests
	@cd $(TEST_PATH) && go test -v -cover -coverpkg $(TEST_PACKAGES) --coverprofile /tmp/$(BINARY).coverage.out ./...
	@sed -E "s/(^.*\/.*:)/.\/src\/\1/g" /tmp/$(BINARY).coverage.out > /tmp/$(BINARY).coverage.fullpath.out

report: # Coverage report (after `make tests`)
	@go tool cover -html=/tmp/$(BINARY).coverage.fullpath.out

	@cd $(TEST_PATH) && go get args && go get -t tests
	@cd $(TEST_PATH) && go test -cpuprofile /tmp/$(BINARY).cpu.prof -memprofile /tmp/$(BINARY).mem.prof -bench .
	@printf "==================================\n"
	@printf "CPU Profiler results...\n"
	@go tool pprof -text /tmp/$(BINARY).cpu.prof
	@printf "==================================\n"
	@printf "Memory Profiler results...\n"
	@go tool pprof -text /tmp/$(BINARY).mem.prof
	@printf "==================================\n"

clean: # Clean all build/test artifacts.
	@printf "Removing build/test artifacts..."
	@if [[ -f $(BIN_DIR)/$(BINARY) ]]; then rm $(BIN_DIR)/$(BINARY); fi
	@if [[ -f $(TEST_PATH)/tests.test ]]; then rm $(TEST_PATH)/tests.test; fi
	@if [[ -d $(BIN_DIR) ]]; then rmdir $(BIN_DIR); fi
	@printf "DONE\n"