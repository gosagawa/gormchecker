VERSION=1.0.0
SOURCE_FILES=$(shell find . -type d -name vendor -prune -o -type d -path ./cmd -prune -o -type f -name '*.go' -print)
GO_LIST=$(shell go list ./... | grep -v /vendor/)
GOPATH=$(shell echo "$$GOPATH")

all: gormchecker

gormchecker: ${SOURCE_FILES} cmd/gormchecker/main.go
	go build -o gormchecker cmd/gormchecker/main.go

test:
	go test -v ./...

