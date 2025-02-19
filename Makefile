BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all
all: audit lint test build

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-permissions-api

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-permissions-api
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-permissions-api

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -race -cover -coverpkg=github.com/ONSdigital/dp-permissions-api/... -component

.PHONY: lint
lint: validate-specification
	golangci-lint run ./...

.PHONY: validate-specification
validate-specification:
	redocly lint swagger.yaml
