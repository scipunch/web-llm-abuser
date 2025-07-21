install:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go mod download

vet: fmt
	go vet ./...
	staticcheck ./...

fmt:
	go fmt ./...

run:
	go run main.go
