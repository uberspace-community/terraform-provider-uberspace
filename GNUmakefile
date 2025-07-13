default: fmt lint install generate

.PHONY: build
build:
	go build -v ./...

.PHONY: install
install: build
	go install -v ./...

.PHONY: install_golangci_lint
install_golangci_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate: fmt
	go tool tfplugindocs generate

.PHONY: fmt
fmt:
	go mod tidy
	golangci-lint run --fix
	terraform fmt -recursive ./examples/
	gofmt -s -w -e .

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=10 ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...
