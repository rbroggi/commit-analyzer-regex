.PHONY: test
test: ## Run Tests into the packages
	@echo "Running tests"
	go test -v -covermode=atomic -coverpkg=./... -coverprofile=cover.out ./...

.PHONY: integration
integration: test
