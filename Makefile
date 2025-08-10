
test:
	go test -v ./...

build-debug:
	mkdir -p build && \
	CGO_ENABLED=1 GO111MODULE=on go build -tags codes -gcflags="all=-N -l" -o build/s3tool ./cmd/s3tool

debug: build-debug
	./build/s3tool

create_test_bucket:
	terraform -chdir=test/deployment/modules/main init && \
	terraform -chdir=test/deployment/modules/main apply -auto-approve

delete_test_bucket:
	terraform -chdir=test/deployment/modules/main init && \
	terraform -chdir=test/deployment/modules/main destroy -auto-approve

generate-images: build-debug
	mkdir -p screens && \
	go run ./cmd/readme
