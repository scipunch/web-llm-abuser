install:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go mod download
	go run github.com/playwright-community/playwright-go/cmd/playwright@v0.5200.0 install --with-deps

vet: fmt
	go vet ./...
	staticcheck ./...

fmt:
	go fmt ./...

run:
	go run cmd/wla/main.go
