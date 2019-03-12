.PHONY: help about prod clean deps deps-go dev docker-build docker-clean docker-dev docker-test v8 targets
.DEFAULT_GOAL := help

LDFLAGS =

V8_VERSION = 6.7.77
v8_dir = ./vendor/github.com/Baldomo/v8-go
v8_include := $(v8_dir)/include
v8_libv8 := $(v8_dir)/libv8

uname := $(shell uname -s)

ifeq ($(OS),Windows_NT)
	target ?= windows
else ifeq ($(uname),Linux)
	target ?= linux
else ifeq ($(uname),Darwin)
	target ?= darwin
endif

help: about targets ## Prints help
about:
	@echo "Fangs' makefile"
targets:
	@echo -e "Make targets:\n"
	@cat $(realpath $(firstword $(MAKEFILE_LIST))) | \
	sed -n -E 's/^([^.][^: ]+)\s*:([^=#]*##\s*(.*[^[:space:]])\s*)$$/    \1: \3/p' | \
	sort -u | \
	expand -t15
	@echo

$(target): ## Builds Fangs for the target OS. Override the variable "target" to build for any other OS
	@export GOOS=$(target)
	@export GOARCH=amd64
	@echo "Building fangs_$(target)"
	@go build -o ./build/fangs_$(target) -ldflags="$(LDFLAGS)" -mod=vendor ./cmd/fangs

prod: LDFLAGS += -s -w
prod: $(current_os) ## Builds Fangs with production flags

clean: ## Cleans up build files
	@echo "Removing /build/"
	@rm -rf build

deps: deps-go v8 ## Download dependencies, tidy and vendor
	@echo "Removing unneeded files from /vendor/"
	@find ./vendor -type f ! \( -name 'modules.txt' -o -name '*.sum' -o -name '*.mod' -o -name '*.rst' -o -name '*.go' -o -name '*.y' -o -name '*.h' -o -name '*.c' -o -name '*.cc' -o -name '*.proto' -o -name '*.tmpl' -o -name '*.s' -o -name '*.pl' -o -name '*.a' \) -exec rm -f {} \;

deps-clean: ## Deletes /vendor and cleans modcache
	@echo "Removing /vendor/"
	@rm -rf vendor

deps-go:
	@echo "Getting"
	@go get $(go list ./... | grep -v /vendor/)
	@echo "Verifying"
	@go mod verify
	@echo "Tidying"
	@go mod tidy
	@echo "Vendoring"
	@go mod vendor

dev: ## Runs Fangs in dev mode (FANGS_DEV set to "true")
	@export FANGS_DEV="true"
	@go run -mod=vendor ./cmd/fangs

docker-build: ## Builds the docker image
	@echo "Building image"
	@docker build -t fangs --force-rm .

docker-clean: ## Removes docker images
	@echo "Killing other instances of Fangs"
	@-docker kill fangs
	@echo "Removing other instances/images of Fangs"
	@-docker rm fangs
	@-docker rmi fangs

docker-dev: docker-clean docker-build ## Deploys Fangs in dev image (FENAGS_DEV set to "true")
	@echo "Running image"
	@docker run --rm -it -p 8080:8080 -p 6060:6060 -e "FANGS_DEV=true" --name fangs fangs

v8: ## Fetches github.com/Baldomo/v8-go and prebuilt V8 headers, if necessary
	@echo "Fetching V8"
	@if [ ! -d $(v8_libv8) ] || [ ! -d $(v8_include) ]; then \
		rm -rf ./vendor/github.com/Baldomo/v8-go; \
		git clone -j 8 https://github.com/Baldomo/v8-go.git ./vendor/github.com/Baldomo/v8-go; \
		docker pull augustoroman/v8-lib:$(V8_VERSION); \
		docker rm v8 || true; \
		docker run --name v8 augustoroman/v8-lib:$(V8_VERSION); \
		cd ./vendor/github.com/Baldomo/v8-go && docker cp v8:/v8/include include/ && docker cp v8:/v8/lib libv8/; \
	fi