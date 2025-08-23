
tests:
	go test -v -cover -coverprofile=coverage.out ./... && \
	go tool cover -html=coverage.out -o coverage.html

build-debug:
	mkdir -p build && \
	CGO_ENABLED=1 GO111MODULE=on go build -tags codes -gcflags="all=-N -l" -o build/s3tool ./cmd/s3tool

debug: build-debug
	./build/s3tool

create_test_bucket:
	tofu -chdir=test/deployment/modules/main init && \
	tofu -chdir=test/deployment/modules/main apply -auto-approve

delete_test_bucket:
	tofu -chdir=test/deployment/modules/main init && \
	tofu -chdir=test/deployment/modules/main destroy -auto-approve

generate-screens:
	mkdir -p screens && \
	go run ./cmd/screens
