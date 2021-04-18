FROM golang:alpine as build
RUN apk --no-cache add tzdata git
WORKDIR /app
ADD . .
RUN GOPRIVATE=git.goasum.de GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o overtime cmd/main.go

FROM scratch as final
COPY --from=build /app/overtime .
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
USER 7000
ENTRYPOINT ["/overtime"]