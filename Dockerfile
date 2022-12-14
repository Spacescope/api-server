FROM golang:1.19.3-bullseye as builder

COPY . /opt
RUN cd /opt && go build -o bin/api-server cmd/api-server/main.go

FROM debian:bullseye
RUN apt update && apt-get install ca-certificates -y
RUN adduser --gecos "Devops Starboard,Github,WorkPhone,HomePhone" --home /app/api-server --disabled-password spacescope
USER spacescope
COPY --from=builder /opt/bin/api-server /app/api-server/api-server

CMD ["--conf", "/app/api-server/service.conf"]
ENTRYPOINT ["/app/api-server/api-server"]
