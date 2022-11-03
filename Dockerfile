FROM golang:1.18.3-bullseye as builder

COPY . /opt
RUN cd /opt && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/data-api-server cmd/api-server/main.go

FROM alpine:3.15.4
RUN mkdir -p /app/api-server-backend
RUN adduser -h /app/api-server-backend -D starboard
USER starboard
COPY --from=builder /opt/bin/data-api-server /app/api-server-backend/data-api-server

CMD ["--conf", "/app/api-server-backend/service.conf"]
ENTRYPOINT ["/app/api-server-backend/data-api-server"]
