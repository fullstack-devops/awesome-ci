PROJECT_PKG = github.com/fullstack-devops/awesome-ci
PKG_LIST = "github.com/fullstack-devops/awesome-ci/cmd/awesome-ci"
BUILD_DIR = ./build

LIST_TARGET_ARCH ?= linux/amd64 linux/arm64 windows/amd64 windows/arm64 darwin/amd64 darwin/arm64
LIST_TARGET_ARCH_COMP ?= linux/amd64 linux/arm64 windows/amd64

LATEST_VERSION ?= "1.0.0"
VERSION ?=$(shell git describe --tags --exact-match 2>/dev/null || echo "dev-pr")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)

# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
# inject build info
LDFLAGS += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}

# all: clean dep test test/cover build
all: clean dep test/cover build

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## dep: fetch dependencies
.PHONY: dep
dep:
	go get -t ./...

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=${BUILD_DIR}/coverage/coverage.out ./...
	go tool cover -html=${BUILD_DIR}/coverage/coverage.out -o ${BUILD_DIR}/coverage/coverage.html

## build: build the application
.PHONY: build
build:
	@echo "[" > ${BUILD_DIR}/package/manifest.json
	@for TARGET in ${LIST_TARGET_ARCH}; do \
		GOOS=$$(echo $$TARGET | cut -d'/' -f1); \
		GOARCH=$$(echo $$TARGET | cut -d'/' -f2); \
		EXT=""; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Building for GOOS=$$GOOS GOARCH=$$GOARCH"; \
		FILENAME="awesome-ci_${VERSION}_$$GOOS-$$GOARCH$$EXT"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package/$$FILENAME ./cmd/awesome-ci; \
		if [ "$$TARGET" = "$$(echo ${LIST_TARGET_ARCH} | awk '{print $$NF}')" ]; then \
			echo "{\"filename\":\"$$FILENAME\", \"path\":\"${BUILD_DIR}/package/$$FILENAME\",\"goos\":\"$$GOOS\",\"goarch\":\"$$GOARCH\"}" >> ${BUILD_DIR}/package/manifest.json; \
		else \
			echo "{\"filename\":\"$$FILENAME\", \"path\":\"${BUILD_DIR}/package/$$FILENAME\",\"goos\":\"$$GOOS\",\"goarch\":\"$$GOARCH\"}," >> ${BUILD_DIR}/package/manifest.json; \
		fi; \
	done
	@echo "]" >> ${BUILD_DIR}/package/manifest.json

## upx: compress binaries
.PHONY: upx
upx:
	@for TARGET in ${LIST_TARGET_ARCH_COMP}; do \
		GOOS=$$(echo $$TARGET | cut -d'/' -f1); \
		GOARCH=$$(echo $$TARGET | cut -d'/' -f2); \
		EXT=""; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Compressing binary for GOOS=$$GOOS GOARCH=$$GOARCH"; \
		upx -5 ./build/package/awesome-ci_${VERSION}_$$GOOS-$$GOARCH$$EXT; \
	done

chglog:
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	cp scripts/release-template.md build/package/release-template.md
	git-chglog ${LATEST_VERSION}.. >> build/package/release-template.md
	git-chglog -o build/package/CHANGELOG.md 1.0.0..

## clean: remove previous builds
.PHONY: clean
clean: ## Remove previous build
	rm -rf ${BUILD_DIR}/coverage/*
	rm -rf ${BUILD_DIR}/docs/*
	rm -rf ${BUILD_DIR}/package/*
	touch ${BUILD_DIR}/docs/.gitkeep ${BUILD_DIR}/coverage/.gitkeep ${BUILD_DIR}/package/.gitkeep
