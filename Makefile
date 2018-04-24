# These will be provided to the target
VERSION := 0.0.1
BUILD := `git rev-parse HEAD`

# Use linker flags to provide version/build settings to the target
LDFLAGS=-tags release -ldflags "-s -w -X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"
ARCH ?= `go env GOHOSTARCH`
GOOS ?= `go env GOOS`

.godeps:
	dep ensure
	touch .godeps

all: agent master rendezvous

test-master:
	@cd master && V=$(V) go test -timeout 3s

test-agents:
	@cd agents/go/cmd && V=$(V) go test -timeout 3s
	@cd agents/go/handler && V=$(V) go test -timeout 3s

package: clean all .godeps
	cd bin && tar -zcvf reverse-shell-$(VERSION)-$(GOOS)-$(ARCH).tar.gz agent master rendezvous

agent: build_dir .godeps
	cd agents/go && go build $(LDFLAGS) -o ../../bin/agent

master: build_dir .godeps
	cd master && go build $(LDFLAGS) -o ../bin/master

rendezvous: build_dir .godeps
	cd rendezvous && go build $(LDFLAGS) -o ../bin/rendezvous

test: all test-master test-agents
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
