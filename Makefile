VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")
LDFLAGS := -ldflags "-X main.Version=${VERSION}"
MODULE = $(shell go list -m)
PACKAGES := $(shell go list ./... | grep -v /vendor/)
GOLINT := ${shell go list -f {{.Target}} golang.org/x/lint/golint}

.PHONY: default
default: help


.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run unit tests
	@echo "mode: count" > coverage-all.out
	@CGO_ENABLED=0 $(foreach pkg,$(PACKAGES), \
		CGO_ENABLED=0 go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg}; \
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test-cover
test-cover: test ## run unit tests and show test coverage information
	go tool cover -html=coverage-all.out


.PHONY: run
run: ## run the API server
	go run ${LDFLAGS} cmd/server/main.go

.PHONY: build
build:  ## build the API server binary
	mkdir -p build
	CGO_ENABLED=0 go build ${LDFLAGS} -a -o build/server $(MODULE)/cmd/server

.PHONY: build-docker
build-docker: ## build the API server as a docker image


.PHONY: clean
clean: ## remove temporary files
	rm -rf server coverage.out coverage-all.out

.PHONY: version
version: ## display the version of the API server
	@echo $(VERSION)

.PHONY: lint
lint: ## run golint on all Go package
	@${GOLINT} $(PACKAGES)

.PHONY: vet
vet: ## run go vet on all Go package
	@CGO_ENABLED=0 go vet $(PACKAGES)

.PHONY: fmt
fmt: ## run go fmt on all Go packages
	@go fmt $(PACKAGES)

.PHONY: start-db
start-db: ## start the database
	@mkdir -p testdata/postgres
	docker run --rm --name postgres -v $(shell pwd)/testdata:/testdata \
		-v $(shell pwd)/testdata/postgres:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pktrade -d -p 5432:5432 postgres

.PHONY: stop-db
stop-db: ## stop the database
	docker stop postgres
