FROM alpine:3.15

RUN mkdir -p /etc/api-server/{bin,conf}
COPY bin/data-api-server /etc/api-server/bin/data-api-server

CMD ["--conf", "/etc/api-server/conf/service.conf"]
ENTRYPOINT ["/etc/api-server/bin/data-api-server"]