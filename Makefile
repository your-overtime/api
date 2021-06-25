linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/overtime_linux cmd/main.go
arm:
	GOOS=linux GOARCH=arm go build -o build/overtime_arm cmd/main.go
mac:
	GOOS=darwin GOARCH=amd64 go build -o build/overtime_darwin cmd/main.go
test:
	go test -v -cover -bench . ./...
test-html:
	go test -coverprofile=coverage.out -bench . ./... && go tool cover -html=coverage.out
