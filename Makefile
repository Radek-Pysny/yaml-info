GO        ?= go
PKGS      ?= ./...
RUN       ?=
TESTFLAGS ?=

ifneq ($(RUN),)
TESTFLAGS += -run '$(RUN)'
endif

.DEFAULT_GOAL := help

.PHONY: help test test-race

help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

test: ## Run unit tests, bypassing the test cache (RUN=FilterRegEx; TESTFLAGS='-extra -flags')
	$(GO) test -count=1 $(TESTFLAGS) $(PKGS)

test-race: ## Run unit tests with the race detector, bypassing the test cache
	$(GO) test -race -count=1 $(TESTFLAGS) $(PKGS)
