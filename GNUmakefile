default: build

build:
	go build -o terraform-provider-americancloud .

install:
	go install .

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

fmt:
	go fmt ./...

lint:
	golangci-lint run ./...

.PHONY: build install test testacc fmt lint
