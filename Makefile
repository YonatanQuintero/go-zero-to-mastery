.PHONY: all test test-race lint sec vuln clean

# Targets
all: lint test

test:
	go test -v ./...

test-race:
	go test -v -race -timeout 30s ./...

lint:
	@which golangci-lint > /dev/null || echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	@which golangci-lint > /dev/null && golangci-lint run ./... || true

vuln:
	@which govulncheck > /dev/null || echo "govulncheck not installed. Run: go install golang.org/x/vuln/cmd/govulncheck@latest"
	@which govulncheck > /dev/null && govulncheck ./... || true

sec:
	@which gosec > /dev/null || echo "gosec not installed. Run: go install github.com/securego/gosec/v2/cmd/gosec@latest"
	@which gosec > /dev/null && gosec -quiet ./... || true

clean:
	rm -f *.out *.pprof trace.out
