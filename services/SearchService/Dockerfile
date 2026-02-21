FROM golang:1.24-alpine AS builder

ENV GRPC_HEALTH_PROBE_VERSION=v0.4.24
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app

COPY services/SearchService/go.mod services/SearchService/go.sum ./services/SearchService/
COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/SearchService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /search-service-app ./services/SearchService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /search-service-app .
CMD ["./search-service-app"]