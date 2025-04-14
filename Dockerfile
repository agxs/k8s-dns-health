FROM golang:1.24 AS builder

WORKDIR /app

COPY . .

RUN make build-linux-amd64

FROM debian:12

COPY --from=builder /app/entrypoint.sh /entrypoint.sh
COPY --from=builder /app/bin/dns-health_linux_amd64 /dns-health

CMD ["/entrypoint.sh"]
