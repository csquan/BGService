all: build


version := "1.0.0"
dockerName := "reg.spb.io/bg/bgservice"

CURRENT_DIR=$(pwd)

PROJ = bgservice
MODULE = "bgservice"

PKG = `go list ./... | grep -v /vendor/`

#print tag into pkg
PKG_TAG = -ldflags "-w -X ${PROJ}/common.gitTag=`git describe --tags` -X ${PROJ}/common.commitNumber=`git rev-parse HEAD` -X ${PROJ}/common.buildTime=`date +%FT%T%z`"

#cross compile
CROSS_COMPILE = CGO_ENABLED=0 GOOS=linux GOARCH=amd64
CILINT := $(shell command -v golangci-lint 2> /dev/null)
GOIMPORTS := $(shell command -v goimports 2> /dev/null)

style:
	! find . -path ./vendor -prune -o -name '*.go' -print | xargs goimports -d -local ${MODULE} | grep '^'

format:
ifndef GOIMPORTS
	$(error "goimports is not available please install goimports")
endif
	find . -path ./vendor -prune -o -name '*.go' -print | xargs goimports -l -local ${MODULE} | xargs goimports -l -local ${MODULE} -w

cilint:
	golangci-lint run

clean:
	rm -rf bin

build:
	go build ${PKG_TAG} -o bin/bgservice main.go
    sudo docker build -t $(dockerName):$(version) .
    sudo docker push $(dockerName):$(version)

test: style cilint
	go test -cover ./...

server: clean 
	${CROSS_COMPILE} go build -o bin/linux-amd64-bgservice ${PKG_TAG} main.go


.PHONY: build clean client
