PROJECT_NAME := "awesome-ci"
PROJECT_PKG = github.com/fullstack-devops/awesome-ci
PKG_LIST = "github.com/fullstack-devops/awesome-ci/cmd/awesome-ci"
BUILD_DIR = ./build

VERSION ?=$(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)

# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
# inject build info
LDFLAGS += -X ${PROJECT_PKG}/internal/app/build.Version=${VERSION} -X ${PROJECT_PKG}/internal/app/build.CommitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/app/build.BuildDate=${BUILD_DATE}

#PLATFORMS := linux/amd64 windows/amd64

.PHONY: docs clean

all: dep awesome-ci coverage

dep:
	go get -t ./...

dep_update:
	go get -t ./...

awesome-ci: dep
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/package ./cmd/awesome-ci

#frontend:
#
test: awesome-ci ## Run unittests
	-go test -short ${PKG_LIST}

race: awesome-ci ## Run data race detector
	-go test -race -short -v ${PKG_LIST}

# this requires ruby with the gem asciidoctor, asciidoctor-pdf and asciidoctor-diagram installed -> gem install asciidoctor-**
# also graphviz is required
docs:
	asciidoctor -b html -r asciidoctor-diagram -d book -D build/docs ./docs/architecture/awesome-ci.adoc
	asciidoctor-pdf -r asciidoctor-diagram -d book -D build/docs ./docs/architecture/awesome-ci.adoc

docspodman:
	podman run --rm -it -v ./:/documents/ docker.io/asciidoctor/docker-asciidoctor asciidoctor -r asciidoctor-diagram -d book -D build/docs ./docs/architecture/awesome-ci.adoc
	podman run --rm -it -v ./:/documents/ docker.io/asciidoctor/docker-asciidoctor asciidoctor-pdf -r asciidoctor-diagram -d book -D build/docs ./docs/architecture/awesome-ci.adoc

coverage:
	-go test -covermode=count -coverprofile "${BUILD_DIR}/coverage/awesome-ci.cov" "git.t3.daimlertruck.com/eMST/skipper/cmd/awesome-ci"
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
	@echo   coverage        - generate test coverage report
	@echo   docs            - generate end user/developer documents
	@echo   docspodman      - generate end user/developer documents with podman
	@echo   clean           - cleanup project direcotories