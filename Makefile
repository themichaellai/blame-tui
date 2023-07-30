src = $(shell find . -type f -name '*.go' -or -name 'go.*' -not -path "./vendor/*")

.PHONY: setup
setup:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.3

bin/blame-tui: $(src)
	go build -o bin/blame-tui ./cmd/blame-tui

.PHONY: lint
lint:
	./bin/golangci-lint run ./...

