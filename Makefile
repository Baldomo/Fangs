.PHONY: dev docker-dev update-deps

LDFLAGS =

current_os =
uname := $(shell uname -s)

ifeq ($(OS),Windows_NT)
	current_os = windows
else ifeq ($(uname),Linux)
	current_os = linux
else ifeq ($(uname),Darwin)
	current_os = darwin
endif

$(current_os):
	@export GOOS=$(current_os)
	@export GOARCH=amd64
	@echo "Building fangs_$(current_os)"
	@go build -o ./build/fangs_$(current_os) -ldflags="$(LDFLAGS)" ./cmd/fangs

prod: LDFLAGS += -s -w
prod: $(current_os)

clean:
	@rm -rf build

dev:
	@export FANGS_DEV="true"
	@go run ./cmd/fangs

docker-dev:
	@echo "Killing other instances of Fangs"
	-docker kill fangs
	@echo "Building image"
	@docker build -t fangs --force-rm . > /dev/null
	@echo "Running image"
	@docker run -it -p 8080:8080 -e "FANGS_DEV=true" --name fangs fangs

update-deps:
	@go get -u ./...
	@go mod verify
	@go mod tidy
	@go mod vendor