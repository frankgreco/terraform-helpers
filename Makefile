ALL_SRC	 = $(shell find . -name "*.go" | grep -v -e vendor)
PACKAGES = $(shell go list ./...)
PASS     = $(shell printf "\033[32mPASS\033[0m")
FAIL     = $(shell printf "\033[31mFAIL\033[0m")
COLORIZE = sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''

.PHONY: fmt
fmt:
	@gofmt -e -s -l -w $(ALL_SRC)

.PHONY: test
test: deps
	@bash -c "set -e; set -o pipefail; go test -v -race $(PACKAGES) | $(COLORIZE)"

.PHONY: deps
deps:
	@go mod download
	@go mod tidy
