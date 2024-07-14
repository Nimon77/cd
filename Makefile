##@ Utilities
help:     ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_\-\.]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

readme:   ## Generate README.md based on GoDoc comments
	go run github.com/posener/goreadme/cmd/goreadme@master -title 'CashDrawer (cd)' -import-path 'github.com/Nimon77/cd' -badge-godoc -generated-notice -types -methods > README.md

##@ Go
build:    ## Build the project
	go build

test:     ## Run Go tests with coverage
	go test -cover -v -json | go run github.com/mfridman/tparse@latest -all