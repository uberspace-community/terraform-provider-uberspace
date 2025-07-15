default: fmt lint install generate

.PHONY: setup
setup:
	go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
	go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest

.PHONY: generate-openapi
generate-openapi: setup
	# replace \u00fc\u00e4\u00f6\u00dc\u00c4\u00d6\u00df with üäöÜÄÖß in openapi.json and save as openapi_fixed.json
	sed 's/\\\\u00fc/ü/g; s/\\\\u00e4/ä/g; s/\\\\u00f6/ö/g; s/\\\\u00dc/Ü/g; s/\\\\u00c4/Ä/g; s/\\\\u00d6/Ö/g; s/\\\\u00df/ß/g' openapi.json > openapi_fixed.json
	tfplugingen-openapi generate \
	  	--config generator_config.yml \
	  	--output provider_code_spec.json \
	  	openapi_fixed.json
	tfplugingen-framework generate all \
        --input provider_code_spec.json \
		--output gen/provider
	go tool ogen --package client --target gen/client --clean openapi_fixed.json

.PHONY: build
build:
	go build -v ./...

.PHONY: install
install: build
	go install -v ./...

.PHONY: install_golangci_lint
install_golangci_lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.2.2

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate: fmt
	go tool tfplugindocs generate

.PHONY: fmt
fmt:
	go mod tidy
	golangci-lint fmt ./...
	terraform fmt -recursive ./examples/
	gofmt -s -w -e .

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=10 ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...
