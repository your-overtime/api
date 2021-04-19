FROM golang:alpine as build
RUN apk --no-cache add tzdata git
FROM scratch as final
COPY build/overtime_linux overtime
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
USER 7000
ENTRYPOINT ["/overtime"]