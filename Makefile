VERSION := 0.0.1
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

BASEDIR := $(shell pwd)
PROJECTNAME := apigw_golang
SOURCEFILES := $(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./src/*" -not -path "./pkg/*" -print)

# init the dir for golang project
# build target for round robin with weight
LOCALENV := /Users/bellke/go
BUILDDIR := ${BASEDIR}/bin
BUILDTARGET := ${BUILDDIR}/apigw_golang
BUILDLINUXTARGET := ${BUILDDIR}/apigw_golang_linux

GOPATH := ${LOCALENV}:${BASEDIR}
GO15VENDOREXPERIMENT := 1
GOENV := CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOPATH=${GOPATH} GO15VENDOREXPERIMENT=${GO15VENDOREXPERIMENT}


LDFLAGS := -X apigw_golang/config.Edition=default -X apigw_golang/config.Version=${VERSION} -X apigw_golang/config.BuildTime=`date +%Y-%m-%dT%T%z`
DEBUGLDFLAGS := -n -X apigw_golang/config.Mode=debug ${LDFLAGS}
RELEASELDFLAGS := -s -w -X apigw_golang/config.Mode=release ${LDFLAGS}

.PHONY: build
build: ${BUILDDIR}
	${GOENV} go build ${BUILDARGS} -i -ldflags "${DEBUGLDFLAGS}" -o ${BUILDTARGET} ${PROJECTNAME}

.PHONY: release
release: ${BUILDDIR}
	${GOENV} go build ${BUILDARGS} -v -ldflags "${RELEASELDFLAGS}" -o ${BUILDTARGET} ${PROJECTNAME}

.PHONY: build-linux
build-linux: ${BUILDDIR}
	GOOS=linux GOARCH=amd64 go build -v -ldflags "${RELEASELDFLAGS}" -o ${BUILDLINUXTARGET} ${PROJECTNAME}

${BUILDDIR}:
	mkdir -p ${BUILDDIR}

.PHONY: test
test:
	cd tests; ${GOENV} go test

.PHONY: codecheck
codecheck:
	${GOENV} find . -type f -name "*.go" -not -path "./vendor/*" -not -path "./src/*" -exec gofmt -d {} \; -exec golint {} \;

.PHONY: clean
clean:
	rm -rf ${BUILDTARGET}
