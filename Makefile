
GOLANGCI_LINT_VERSION ?= v2.11.4

.PHONY: test test-cover tests build build-debug debug fmt vet lint lint-install docs-commands create_test_bucket delete_test_bucket generate-screens

test:
	go test -v ./...

test-cover:
	go test -v -cover -coverprofile=coverage.out ./... && \
	go tool cover -html=coverage.out -o coverage.html

tests: test-cover

build:
	mkdir -p build && \
	CGO_ENABLED=1 GO111MODULE=on go build -o build/s3tool ./cmd/s3tool

build-debug:
	mkdir -p build && \
	CGO_ENABLED=1 GO111MODULE=on go build -tags codes -gcflags="all=-N -l" -o build/s3tool ./cmd/s3tool

debug: build-debug
	./build/s3tool

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run --timeout=5m ./...

lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

docs-commands:
	mkdir -p docs && { \
		echo "# Command Reference"; \
		echo; \
		echo "Generated from Cobra help output."; \
		echo; \
		echo "## s3tool --help"; \
		echo "\`\`\`text"; \
		go run ./cmd/s3tool --help; \
		echo "\`\`\`"; \
		echo; \
		echo "## s3tool completion --help"; \
		echo "\`\`\`text"; \
		go run ./cmd/s3tool completion --help; \
		echo "\`\`\`"; \
	} > docs/COMMANDS.md

create_test_bucket:
	tofu -chdir=test/deployment/modules/main init && \
	tofu -chdir=test/deployment/modules/main apply -auto-approve

delete_test_bucket:
	tofu -chdir=test/deployment/modules/main init && \
	tofu -chdir=test/deployment/modules/main destroy -auto-approve

generate-screens:
	mkdir -p screens && \
	go run ./cmd/screens
