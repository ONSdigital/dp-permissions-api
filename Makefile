BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

JAVA_SDK_DIR="./sdk-java"

.PHONY: all
all: audit lint test build

.PHONY: audit
audit: audit-go audit-java

.PHONY: audit-go
audit-go:
	dis-vulncheck

.PHONY: audit-java
audit-java: 
	mvn -f $(JAVA_SDK_DIR) ossindex:audit

.PHONY: build
build: build-go build-java

.PHONY: build-go
build-go:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-permissions-api

.PHONY: build-java
build-java:	
	mvn -f $(JAVA_SDK_DIR) clean package -Dmaven.test.skip -Dossindex.skip=true

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-permissions-api
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-permissions-api

.PHONY: test
test: test-go test-java

.PHONY: test-go
test-go: ## Runs unit tests including checks for race conditions and returns coverage
	go test -race -cover ./...

.PHONY: test-java
test-java:
	mvn -f $(JAVA_SDK_DIR) -Dossindex.skip=true test

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -race -cover -coverpkg=github.com/ONSdigital/dp-permissions-api/... -component

.PHONY: lint
lint: lint-go lint-java validate-specification

.PHONY: lint-go
lint-go: ## Used in ci to run linters against Go code
	golangci-lint run ./...

.PHONY: lint-java
lint-java:
	mvn -f $(JAVA_SDK_DIR) clean checkstyle:check test-compile

.PHONY: validate-specification
validate-specification:
	redocly lint swagger.yaml
