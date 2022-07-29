PWD := $(shell pwd)
GOPATH := $(shell go env GOPATH)
LDFLAGS := $(shell go run buildscripts/gen-ldflags.go)

GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

VERSION ?= $(shell git describe --tags)
TAG ?= "minio/minio:$(VERSION)"

all: build

checks: ## check dependencies
	@echo "Checking dependencies"
	@(env bash $(PWD)/buildscripts/checkdeps.sh)

help: ## print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

getdeps: ## fetch necessary dependencies
	@mkdir -p ${GOPATH}/bin
	@echo "Installing golangci-lint" && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.45.2
	@echo "Installing msgp" && go install -v github.com/tinylib/msgp@v1.1.7-0.20211026165309-e818a1881b0e
	@echo "Installing stringer" && go install -v golang.org/x/tools/cmd/stringer@latest

crosscompile: ## cross compile minio
	@(env bash $(PWD)/buildscripts/cross-compile.sh)

verifiers: getdeps lint check-gen

check-gen: ## check for updated autogenerated files
	@go generate ./... >/dev/null
	@(! git diff --name-only | grep '_gen.go$$') || (echo "Non-committed changes in auto-generated code is detected, please commit them to proceed." && false)

lint: ## runs golangci-lint suite of linters
	@echo "Running $@ check"
	@${GOPATH}/bin/golangci-lint run --build-tags kqueue --timeout=10m --config ./.golangci.yml

check: test
test: verifiers build ## builds minio, runs linters, tests
	@echo "Running unit tests"
	@MINIO_API_REQUESTS_MAX=10000 CGO_ENABLED=0 go test -tags kqueue ./...

test-decom: install
	@echo "Running minio decom tests"
	@env bash $(PWD)/docs/distributed/decom.sh
	@env bash $(PWD)/docs/distributed/decom-encrypted.sh
	@env bash $(PWD)/docs/distributed/decom-encrypted-sse-s3.sh
	@env bash $(PWD)/docs/distributed/decom-compressed-sse-s3.sh

test-upgrade: build
	@echo "Running minio upgrade tests"
	@(env bash $(PWD)/buildscripts/minio-upgrade.sh)

test-race: verifiers build ## builds minio, runs linters, tests (race)
	@echo "Running unit tests under -race"
	@(env bash $(PWD)/buildscripts/race.sh)

test-iam: build ## verify IAM (external IDP, etcd backends)
	@echo "Running tests for IAM (external IDP, etcd backends)"
	@MINIO_API_REQUESTS_MAX=10000 CGO_ENABLED=0 go test -tags kqueue -v -run TestIAM* ./cmd
	@echo "Running tests for IAM (external IDP, etcd backends) with -race"
	@MINIO_API_REQUESTS_MAX=10000 GORACE=history_size=7 CGO_ENABLED=1 go test -race -tags kqueue -v -run TestIAM* ./cmd

test-replication: install ## verify multi site replication
	@echo "Running tests for replicating three sites"
	@(env bash $(PWD)/docs/bucket/replication/setup_3site_replication.sh)
	@(env bash $(PWD)/docs/bucket/replication/setup_2site_existing_replication.sh)

test-site-replication-ldap: install ## verify automatic site replication
	@echo "Running tests for automatic site replication of IAM (with LDAP)"
	@(env bash $(PWD)/docs/site-replication/run-multi-site-ldap.sh)

test-site-replication-oidc: install ## verify automatic site replication
	@echo "Running tests for automatic site replication of IAM (with OIDC)"
	@(env bash $(PWD)/docs/site-replication/run-multi-site-oidc.sh)

test-site-replication-minio: install ## verify automatic site replication
	@echo "Running tests for automatic site replication of IAM (with MinIO IDP)"
	@(env bash $(PWD)/docs/site-replication/run-multi-site-minio-idp.sh)

verify: ## verify minio various setups
	@echo "Verifying build with race"
	@GORACE=history_size=7 CGO_ENABLED=1 go build -race -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null
	@(env bash $(PWD)/buildscripts/verify-build.sh)

verify-healing: ## verify healing and replacing disks with minio binary
	@echo "Verify healing build with race"
	@GORACE=history_size=7 CGO_ENABLED=1 go build -race -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null
	@(env bash $(PWD)/buildscripts/verify-healing.sh)
	@(env bash $(PWD)/buildscripts/unaligned-healing.sh)

verify-healing-with-root-disks: ## verify healing root disks
	@echo "Verify healing with root disks"
	@GORACE=history_size=7 CGO_ENABLED=1 go build -race -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null
	@(env bash $(PWD)/buildscripts/verify-healing-with-root-disks.sh)

verify-healing-with-rewrite: ## verify healing to rewrite old xl.meta -> new xl.meta
	@echo "Verify healing with rewrite"
	@GORACE=history_size=7 CGO_ENABLED=1 go build -race -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null
	@(env bash $(PWD)/buildscripts/rewrite-old-new.sh)

verify-healing-inconsistent-versions: ## verify resolving inconsistent versions
	@echo "Verify resolving inconsistent versions build with race"
	@GORACE=history_size=7 CGO_ENABLED=1 go build -race -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null
	@(env bash $(PWD)/buildscripts/resolve-right-versions.sh)

build: checks ## builds minio to $(PWD)
	@echo "Building minio binary to './minio'"
	@CGO_ENABLED=0 go build -tags kqueue -trimpath --ldflags "$(LDFLAGS)" -o $(PWD)/minio 1>/dev/null

hotfix-vars:
	$(eval LDFLAGS := $(shell MINIO_RELEASE="RELEASE" MINIO_HOTFIX="hotfix.$(shell git rev-parse --short HEAD)" go run buildscripts/gen-ldflags.go $(shell git describe --tags --abbrev=0 | \
    sed 's#RELEASE\.\([0-9]\+\)-\([0-9]\+\)-\([0-9]\+\)T\([0-9]\+\)-\([0-9]\+\)-\([0-9]\+\)Z#\1-\2-\3T\4:\5:\6Z#')))
	$(eval VERSION := $(shell git describe --tags --abbrev=0).hotfix.$(shell git rev-parse --short HEAD))
	$(eval TAG := "minio/minio:$(VERSION)")

hotfix: hotfix-vars install ## builds minio binary with hotfix tags
	@mv -f ./minio ./minio.$(VERSION)
	@minisign -qQSm ./minio.$(VERSION) -s "${CRED_DIR}/minisign.key" < "${CRED_DIR}/minisign-passphrase"
	@sha256sum < ./minio.$(VERSION) | sed 's, -,minio.$(VERSION),g' > minio.$(VERSION).sha256sum

hotfix-push: hotfix
	@scp -q -r minio.$(VERSION)* minio@dl-0.minio.io:~/releases/server/minio/hotfixes/linux-amd64/archive/
	@scp -q -r minio.$(VERSION)* minio@dl-1.minio.io:~/releases/server/minio/hotfixes/linux-amd64/archive/
	@echo "Published new hotfix binaries at https://dl.min.io/server/minio/hotfixes/linux-amd64/archive/minio.$(VERSION)"

docker-hotfix-push: docker-hotfix
	@docker push -q $(TAG) && echo "Published new container $(TAG)"

docker-hotfix: hotfix-push checks ## builds minio docker container with hotfix tags
	@echo "Building minio docker image '$(TAG)'"
	@docker build -q --no-cache -t $(TAG) --build-arg RELEASE=$(VERSION) . -f Dockerfile.hotfix

docker: build checks ## builds minio docker container
	@echo "Building minio docker image '$(TAG)'"
	@docker build -q --no-cache -t $(TAG) . -f Dockerfile

install: build ## builds minio and installs it to $GOPATH/bin.
	@echo "Installing minio binary to '$(GOPATH)/bin/minio'"
	@mkdir -p $(GOPATH)/bin && cp -f $(PWD)/minio $(GOPATH)/bin/minio
	@echo "Installation successful. To learn more, try \"minio --help\"."

clean: ## cleanup all generated assets
	@echo "Cleaning up all the generated files"
	@find . -name '*.test' | xargs rm -fv
	@find . -name '*~' | xargs rm -fv
	@find . -name '.#*#' | xargs rm -fv
	@find . -name '#*#' | xargs rm -fv
	@rm -rvf minio
	@rm -rvf build
	@rm -rvf release
	@rm -rvf .verify*
