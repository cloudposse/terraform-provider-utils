TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=cloudposse
NAME=utils
BINARY=terraform-provider-${NAME}
VERSION=9999.99.99
#GOOS=darwin
GOOS=linux
GOARCH=amd64
SHELL := /bin/bash

# List of targets the `readme` target should call before generating the readme
export README_DEPS ?= docs/targets.md docs/terraform.md

-include $(shell curl -sSL -o .build-harness "https://cloudposse.tools/build-harness"; echo .build-harness)

get:
	go get

build: get
	env GOOS=${GOOS} GOARCH=${GOARCH} go build

deps:
	go mod download

docs:
	go generate

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}

# Lint terraform code
lint:
	$(SELF) terraform/install terraform/get-modules terraform/get-plugins terraform/lint terraform/validate

# Run acceptance tests
testacc: install
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 3m

tfplugindocs:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
	tfplugindocs generate
