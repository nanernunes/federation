GOCMD:=$(shell which go)
GOSWAG:=$(shell which swag)

test:
	@$(GOCMD) test -v ./...

cover:
	@$(GOCMD) test -v ./... -coverprofile=coverage.txt -covermode=atomic

deps:
	@$(GOCMD) mod download

build:
	@$(GOCMD) build -o federation federation.go

run:
	@$(GOCMD) run federation.go
