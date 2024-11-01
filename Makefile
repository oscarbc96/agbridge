ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), see https://golangci-lint.run/usage/install/#local-installation")
endif

ifeq (, $(shell which goreleaser))
$(warning "could not find goreleaser in $(PATH), see https://goreleaser.com/install/")
endif

ifeq (, $(shell which gotestsum))
$(warning "could not find gotestsum in $(PATH), see https://github.com/gotestyourself/gotestsum#install")
endif

ifeq (, $(shell which gofumpt))
$(warning "could not find gofumpt in $(PATH), see https://github.com/mvdan/gofumpt")
endif

.PHONY: all format lint test snapshot install freeze freeze-upgrade generate

default: all

all: install generate lint test snapshot

format:
	$(info ******************** checking formatting ********************)
	gofumpt -w -extra .

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test:
	$(info ******************** running tests ********************)
	gotestsum -- -coverprofile=cover.out ./...

snapshot:
	$(info ******************** building bin: snapshot ********************)
	goreleaser build --clean --snapshot --single-target -o "dist/agbridge"
	@echo "Bin available at: dist/agbridge"

install:
	$(info ******************** downloading dependencies ********************)
	go mod download

freeze:
	$(info ******************** freeze dependencies ********************)
	go mod tidy && go mod verify

freeze-upgrade:
	$(info ******************** upgrade dependencies ********************)
	go get -u ./... && go mod tidy && go mod verify

generate:
	$(info ******************** generating support files ********************)
	go generate ./...
