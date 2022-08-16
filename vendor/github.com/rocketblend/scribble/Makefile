.DEFAULT_GOAL := check
check: vet race fmt ## Check project
	
vet: ## Vet the files
	@go vet

fmt: ## Style check the files
	@gofmt -s -w .

test: ## Run tests
	@go test -short

race: ## Run tests with data race detector
	@go test -race

benchmark: ## Run benchmarks
	@go test -run="-" -bench=".*"

dep:
	@go mod download
	@go mod tidy
