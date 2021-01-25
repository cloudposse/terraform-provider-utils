default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 TF_ACC_PROVIDER_HOST=registry.terrform.io TF_ACC_PROVIDER_NAMESPACE=cloudposse go test ./... -v $(TESTARGS) -timeout 120m
