linux:
	GOPRIVATE=git.goasum.de GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/tracer_linux cmd/main.go
arm:
	GOPRIVATE=git.goasum.de GOOS=linux GOARCH=arm go build -o build/tracer_arm cmd/main.go
mac:
	GOPRIVATE=git.goasum.de GOOS=darwin GOARCH=amd64 go build -o build/tracer_darwin cmd/main.go
test:
	GOPRIVATE=git.goasum.de go test -v -cover -bench . ./...
test-html:
	GOPRIVATE=git.goasum.de go test -coverprofile=coverage.out -bench . ./... && go tool cover -html=coverage.out
