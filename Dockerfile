FROM golang:1.18.3-bullseye as builder

COPY . /opt
RUN cd /opt && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-server cmd/api-server/main.go

FROM alpine:3.15.4
RUN mkdir -p /app/api-server
RUN adduser -h /app/api-server -D starboard
USER starboard
COPY --from=builder /opt/bin/api-server /app/api-server/api-server

CMD ["--conf", "/app/api-server/service.conf"]
ENTRYPOINT ["/app/api-server/api-server"]
