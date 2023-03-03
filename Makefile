PROJECT_NAME := "awesome-ci"
PROJECT_PKG = github.com/fullstack-devops/awesome-ci
PKG_LIST = "github.com/fullstack-devops/awesome-ci/cmd/awesome-ci"
BUILD_DIR = ./build

LATEST_VERSION ?= "1.0.0"
VERSION ?=$(shell git describe --tags --exact-match 2>/dev/null || echo "dev-pr")
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)


# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
# inject build info
LDFLAGS += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}

.PHONY: docs clean

all: dep awesome-ci
# coverage

dep:
	go get -t ./...

dep_update:
	go get -t ./...

awesome-ci: dep
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package/awesome-ci_${VERSION}_linux-amd64 ./cmd/awesome-ci
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package/awesome-ci_${VERSION}_linux-arm64 ./cmd/awesome-ci
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package/awesome-ci_${VERSION}_windows-amd64.exe ./cmd/awesome-ci
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package/awesome-ci_${VERSION}_windows-arm64.exe ./cmd/awesome-ci

test: ## Run unittests
	-go test -short -v ./internal/pkg/...

race: ## Run data race detector
	-go test -race -short -v ${PKG_LIST}

chglog:
	go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	cp scripts/release-template.md build/package/release-template.md
	git-chglog ${LATEST_VERSION}.. >> build/package/release-template.md
	git-chglog -o build/package/CHANGELOG.md 1.0.0..

coverage:
	-go test -covermode=count -coverprofile "${BUILD_DIR}/coverage/awesome-ci.cov" "github.com/fullstack-devops/awesome-ci/cmd/awesome-ci"
	echo mode: count > "${BUILD_DIR}/coverage/coverage.cov"
	tail -n +2 "${BUILD_DIR}/coverage/awesome-ci.cov" >> "${BUILD_DIR}/coverage/coverage.cov"
	go tool cover -html="${BUILD_DIR}/coverage/coverage.cov" -o "${BUILD_DIR}/coverage/coverage.html"


clean: ## Remove previous build
	rm -rf ${BUILD_DIR}/docs/*
	rm -rf ${BUILD_DIR}/coverage/*
	rm -rf ${BUILD_DIR}/package/*
	touch ${BUILD_DIR}/docs/.keep ${BUILD_DIR}/coverage/.keep ${BUILD_DIR}/package/.keep


help:
	@echo Available targets are:
	@echo   all             - build all
	@echo   dep             - fetch dependencies
	@echo   dep_update      - update dependencies
	@echo   awesome-ci      - build awesome-ci
	@echo   test            - run tests
	@echo   race            - run race condition tests
	@echo   chglog          - install and create changelog with chglog
	@echo   coverage        - generate test coverage report
	@echo   clean           - cleanup project direcotories