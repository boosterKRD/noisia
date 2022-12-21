APPNAME = noisia

TAG=$(shell git describe --tags |cut -d- -f1)
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS = -a -installsuffix cgo -ldflags "-X main.appName=${APPNAME} -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"

.PHONY: help clean dep lint test build build-docker

.DEFAULT_GOAL := help

help: ## Display this help screen
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Clean
	rm -f ./bin/${APPNAME}
	rmdir ./bin

dep: ## Get the dependencies
	go mod download

lint: ## Lint the source files
#	golangci-lint run --timeout 5m -E golint -e '(method|func) [a-zA-Z]+ should be [a-zA-Z]+'
#	https://staticcheck.io/docs/getting-started/
	staticcheck ./...

test: dep ## Run tests
	go test -race -timeout 300s -coverprofile=.test_coverage.txt ./... && \
    	go tool cover -func=.test_coverage.txt | tail -n1 | awk '{print "Total test coverage: " $$3}'
	@rm .test_coverage.txt

build: dep ## Build
	mkdir -p ./bin
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${APPNAME} ./cmd

docker-build: ## Build docker image
	docker build --build-arg noisia_my_msg="It's my TEXT" -t boosterkrd/${APPNAME}:${TAG} .
	docker image prune --force --filter label=stage=intermediate
	docker tag boosterkrd/${APPNAME}:${TAG} boosterkrd/${APPNAME}:latest

docker-push: ## Push docker image to the registry
	docker push boosterkrd/${APPNAME}:${TAG}
	docker push boosterkrd/${APPNAME}:latest
