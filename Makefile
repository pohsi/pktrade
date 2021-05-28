VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || echo "1.0.0")

.PHONY: default
default: help


.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test
test: ## run unit tests

.PHONY: run
run: ## run the API server

.PHONY: build
build:  ## build the API server binary

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
	@golint $(PACKAGES)

.PHONY: fmt
fmt: ## run "go fmt" on all Go packages
	@go fmt $(PACKAGES)

.PHONY: start-db
start-db: ## start the database
	@mkdir -p data/postgres
	docker run --rm --name postgres -v $(shell pwd)/data:/data \
		-v $(shell pwd)/data/postgres:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pktrade -d -p 8001:8001 postgres

.PHONY: stop-db
stop-db: ## stop the database
	docker stop postgres
