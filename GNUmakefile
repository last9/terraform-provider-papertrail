TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=transloadit
GOPATH?=$(HOME)/go

PKG_OS?= windows darwin linux 
PKG_ARCH?= amd64
BASE_PATH?= $(shell pwd)
BUILD_PATH?= $(BASE_PATH)/build
PROVIDER := "terraform-provider-transloadit"
VERSION ?= v0.0.0

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=10

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout=120m -parallel=10

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

release: vet test
	@for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			mkdir -p $(BUILD_PATH) && \
			cd $(BASE_PATH) && \
			if [ "$${os}" = "windows" ]; then OS_EXT=".exe"; echo "windows ext"; else OS_EXT=""; fi && \
			echo "$${os} $${OS_EXT}" && \
			rm -f $(BUILD_PATH)/$(PROVIDER)_$(VERSION)$${OS_EXT} && \
			cgo_enabled=0 GOOS=$${os} GOARCH=$${arch} go build -o $(BUILD_PATH)/$(PROVIDER)_$(VERSION)$${OS_EXT} . && \
			cd $(BUILD_PATH) && \
			rm -f $(BUILD_PATH)/$(PROVIDER)_$${os}_$${arch}.tar.gz && \
			tar -cvzf $(BUILD_PATH)/$(PROVIDER)_$${os}_$${arch}.tar.gz $(PROVIDER)_$(VERSION)$${OS_EXT} && \
			rm -f $(PROVIDER)_$(VERSION)$${OS_EXT}; \
		done; \
	done;


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	cd $(GOPATH)/src/$(WEBSITE_REPO) && git submodule init ext/providers/$(PKG_NAME) && git submodule update
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), getting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	cd $(GOPATH)/src/$(WEBSITE_REPO) && git submodule init ext/providers/$(PKG_NAME) && git submodule update
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc vet fmt fmtcheck errcheck test-compile website website-test
