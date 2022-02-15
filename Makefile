GOCMD:=$(shell which go)
GOSWAG:=$(shell which swag)

deps:
	@$(GOCMD) mod download

build:
	@$(GOCMD) build -o federation federation.go

run:
	@$(GOCMD) run federation.go
