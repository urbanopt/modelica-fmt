.PHONY: help all build test clean generate-parser update-testdata release-check release

.DEFAULT_GOAL := help

GO ?= go
BINARY ?= modelica-fmt
REMOTE ?= origin
VERSION ?=
TAG ?= $(if $(VERSION),$(if $(filter v%,$(VERSION)),$(VERSION),v$(VERSION)))
MESSAGE ?= Release $(TAG)
SHA ?= HEAD
VERSION_ID := $(patsubst v%,%,$(TAG))
COMMIT := $(shell git rev-parse --short=12 $(SHA) 2>/dev/null)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
RELEASE_LDFLAGS := -X main.version=$(VERSION_ID) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=make

help:
	@echo "Available targets:"
	@echo "  all              build and test"
	@echo "  build            build $(BINARY)"
	@echo "  test             run Go tests"
	@echo "  clean            remove build artifacts"
	@echo "  generate-parser  regenerate parser from grammar"
	@echo "  update-testdata  regenerate formatter test data"
	@echo "  release-check    verify release build metadata (VERSION=x.y.z)"
	@echo "  release          check, tag, and push a release (VERSION=x.y.z)"

all: build test

build:
	$(GO) build -o $(BINARY) .

test:
	$(GO) test ./...

clean:
	rm -f $(BINARY)

generate-parser:
	./generate_parser.sh

update-testdata:
	./update_testdata.sh

release-check: test
	@if [ -z "$(TAG)" ]; then \
		echo "error: set VERSION=x.y.z or TAG=vx.y.z"; \
		exit 2; \
	fi
	@printf '%s\n' "$(TAG)" | grep -Eq '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z.-]+)?$$' || { \
		echo "error: TAG must be a semantic version tag such as v1.2.3 or v1.2.3-rc.1"; \
		exit 2; \
	}
	@if git rev-parse -q --verify "refs/tags/$(TAG)" >/dev/null; then \
		echo "error: tag $(TAG) already exists"; \
		exit 2; \
	fi
	@git rev-parse --verify "$(SHA)^{commit}" >/dev/null || { \
		echo "error: SHA must point to a commit"; \
		exit 2; \
	}
	@command -v goreleaser >/dev/null 2>&1 || { \
		echo "error: goreleaser is required for release-check"; \
		exit 2; \
	}
	$(GO) build -ldflags "$(RELEASE_LDFLAGS)" -o $(BINARY) .
	@version_output=$$(./$(BINARY) -v); \
	printf '%s\n' "$$version_output"; \
	printf '%s\n' "$$version_output" | grep -F "modelicafmt v$(VERSION_ID)" >/dev/null || { \
		echo "error: expected $(BINARY) -v to report version $(VERSION_ID)"; \
		exit 2; \
	}
	goreleaser check
	goreleaser release --snapshot --clean

release: release-check
	git tag -a "$(TAG)" -m "$(MESSAGE)" "$(SHA)"
	git push "$(REMOTE)" "$(TAG)"
