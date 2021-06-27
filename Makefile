build:
	GOOS=linux GOARCH=amd64 go build -o bin/nmg-linux-amd64 cmd/nmg/nmg.go
	GOOS=darwin GOARCH=amd64 go build -o bin/nmg-mac-amd64 cmd/nmg/nmg.go
	GOOS=windows GOARCH=amd64 go build -o bin/nmg-windows-amd64.exe cmd/nmg/nmg.go

	GOOS=linux GOARCH=arm64 go build -o bin/nmg-linux-arm64 cmd/nmg/nmg.go
	GOOS=darwin GOARCH=arm64 go build -o bin/nmg-mac-arm64 cmd/nmg/nmg.go

run:
	go run cmd/nmb/nmg.go

test:
	go test ./...
