linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/overtime_linux main.go
arm:
	GOOS=linux GOARCH=arm go build -o build/overtime_arm main.go
mac:
	GOOS=darwin GOARCH=amd64 go build -o build/overtime_darwin main.go
test:
	go test -v -cover -bench . ./...
test-html:
	go test -coverprofile=coverage.out -bench . ./... && go tool cover -html=coverage.out
swagger:
	${GOPATH}/bin/swag init --parseDependency --parseInternal -p pascalcase