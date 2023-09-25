default: test

test: ## Run all tests
	$Q go test -v ./... -coverprofile=coverage.out -covermode=atomic