default: testacc

.PHONY: build
build:
	go build

.PHONY: docs
docs:
	go generate

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
