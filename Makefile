default: test

tests: ## Run all tests
	$Q go test -v ./... -coverprofile=coverage.out -cover -coverpkg ./...