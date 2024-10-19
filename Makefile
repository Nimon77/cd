##@ Utilities
help:    ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_\-\.]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Go
build:   ## Build the project
	go build

test:    ## Run Go tests with coverage
	go test -cover -v -json | go run github.com/mfridman/tparse@latest -all

cover:   ## Generate HTML coverage report
	go test -v -coverprofile cover.out -json | go run github.com/mfridman/tparse@latest -all
	go tool cover -html cover.out -o cover.html
	@echo "Coverage report in cover.html"
	go tool cover -html cover.out
