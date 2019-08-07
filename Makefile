# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-tags release -ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
ARCH ?= `go env GOHOSTARCH`
GOOS ?= `go env GOOS`

all: agent client rendezvous

test-client:
	@cd client && V=$(V) go test -timeout 3s

test-agent:
	@cd agent/cmd && V=$(V) go test -timeout 3s
	@cd agent/handler && V=$(V) go test -timeout 3s

package: clean all docs
	cd bin && tar -zcvf reverse-shell-$(VERSION)-$(GOOS)-$(ARCH).tar.gz reverse-shell-agent reverse-shell-client reverse-shell-rendezvous

agent: build_dir
	cd agent && go build $(LDFLAGS) -o ../bin/reverse-shell-agent

client: build_dir
	cd client && go build $(LDFLAGS) -o ../bin/reverse-shell-client

rendezvous: build_dir
	cd rendezvous && go build $(LDFLAGS) -o ../bin/reverse-shell-rendezvous

doc-generator: build_dir
	cd docs && go build $(LDFLAGS) -o ../bin/doc-generator

docs: doc-generator
	cd docs && rm -rf agent/* client/* rendezvous/*
	./bin/doc-generator

test: all test-client test-agent
	@true

build_dir:
	mkdir -p bin

clean:
	rm -rf bin/*

.PHONY: cover
cover:
	@rm -rf coverage.txt
	@for d in `go list ./...`; do \
		t=$$(date +%s); \
		go test -coverprofile=cover.out -covermode=atomic $$d || exit 1; \
		echo "Coverage test $$d took $$(($$(date +%s)-t)) seconds"; \
		if [ -f cover.out ]; then \
			cat cover.out >> coverage.txt; \
			rm cover.out; \
		fi; \
	done
	@echo "Uploading coverage results..."
	@curl -s https://codecov.io/bash | bash
