.PHONY: help all build test clean generate-parser update-testdata changelog-check release-check release

.DEFAULT_GOAL := help

GO ?= go
BINARY ?= modelica-fmt
REMOTE ?= origin
CHANGELOG ?= CHANGELOG.md
VERSION ?=
TAG ?= $(if $(VERSION),$(if $(filter v%,$(VERSION)),$(VERSION),v$(VERSION)))
MESSAGE ?= Release $(TAG)
SHA ?= HEAD
VERSION_ID := $(patsubst v%,%,$(TAG))
COMMIT := $(shell git rev-parse --short=12 $(SHA) 2>/dev/null)
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
RELEASE_DATE := $(shell date -u +%Y-%m-%d)
RELEASE_LDFLAGS := -X main.version=$(VERSION_ID) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=make

help:
	@echo "Available targets:"
	@echo "  all              build and test"
	@echo "  build            build $(BINARY)"
	@echo "  test             run Go tests"
	@echo "  clean            remove build artifacts"
	@echo "  generate-parser  regenerate parser from grammar"
	@echo "  update-testdata  regenerate formatter test data"
	@echo "  changelog-check  verify $(CHANGELOG) has unreleased notes"
	@echo "  release-check    verify release build metadata (VERSION=x.y.z)"
	@echo "  release          check, tag, and push a release (VERSION=x.y.z)"
	@echo ""
	@echo "$(CHANGELOG) is the source of truth for release notes: entries added"
	@echo "under '## [Unreleased]' are moved into a dated '## [vX.Y.Z]' section by"
	@echo "'make release' and published verbatim as the GitHub release body."

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

changelog-check:
	@notes="$$(scripts/changelog.sh extract Unreleased $(CHANGELOG))"; \
	if [ -z "$$(printf '%s' "$$notes" | tr -d '[:space:]')" ]; then \
		echo "error: $(CHANGELOG) has no entries under '## [Unreleased]'; add release notes before releasing"; \
		exit 2; \
	fi

release-check: test changelog-check
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
	@notes_file="$$(mktemp)"; \
	trap 'rm -f "$$notes_file"' EXIT; \
	scripts/changelog.sh extract Unreleased $(CHANGELOG) > "$$notes_file"; \
	goreleaser check; \
	goreleaser release --snapshot --clean --release-notes "$$notes_file"

release: release-check
	@if [ "$$(git rev-parse $(SHA))" != "$$(git rev-parse HEAD)" ]; then \
		echo "error: automatic $(CHANGELOG) finalization requires SHA=HEAD (got $(SHA));"; \
		echo "       finalize $(CHANGELOG) for $(TAG) manually and retry, or release from HEAD"; \
		exit 2; \
	fi
	scripts/changelog.sh finalize "$(TAG)" "$(RELEASE_DATE)" $(CHANGELOG)
	git add $(CHANGELOG)
	git commit -m "docs: update $(CHANGELOG) for $(TAG)"
	git tag -a "$(TAG)" -m "$(MESSAGE)" HEAD
	git push "$(REMOTE)" HEAD
	git push "$(REMOTE)" "$(TAG)"
