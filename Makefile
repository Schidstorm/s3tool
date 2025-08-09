

# create function
define build_debug
	$(shell mkdir -p debug-build && \
		    CGO_ENABLED=1 GO111MODULE=on go build -tags codes -gcflags="all=-N -l" -o debug-build/$(1) ./cmd/$(1) && \
			./debug-build/$(1))
endef

test:
	go test -v ./...

debug:
	$(call build_debug,s3tool)

create_test_bucket:
	terraform -chdir=test/deployment/modules/main init && \
	terraform -chdir=test/deployment/modules/main apply -auto-approve

delete_test_bucket:
	terraform -chdir=test/deployment/modules/main init && \
	terraform -chdir=test/deployment/modules/main destroy -auto-approve

