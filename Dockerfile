FROM golang:1.14.2-alpine as base

RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -installsuffix cgo -o app cmd/main.go

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.1 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM scratch

COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /build/app /bin/
COPY --from=base /bin/grpc_health_probe /bin/grpc_health_probe

ENTRYPOINT [ "/bin/app" ]