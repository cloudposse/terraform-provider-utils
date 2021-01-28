TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=registry.terraform.io
NAMESPACE=cloudposse
NAME=utils
BINARY=terraform-provider-${NAME}
VERSION=9999.99.99
OS_ARCH=darwin_amd64

default: testacc

build:
	go build

deps:
	go mod download

docs:
	go generate

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

# Run acceptance tests
testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
