code-check:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./clients/...
	staticcheck ./dto/...
	staticcheck ./query/...

check-security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)