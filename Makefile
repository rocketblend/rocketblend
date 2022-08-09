tidy:
	go mod tidy

vendor:
	go mod vendor

run:
	go run ./cmd/cli

launch:
	go run ./cmd/launcher