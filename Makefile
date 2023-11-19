check-security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

code-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./clients/...
	staticcheck ./dto/...
	staticcheck ./query/...
	make check-security
	make tests

code-clean:
	make imports
	make format

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)

tests:
	go test ./...
