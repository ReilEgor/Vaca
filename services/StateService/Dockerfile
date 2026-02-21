FROM golang:1.24-alpine AS builder

ENV GRPC_HEALTH_PROBE_VERSION=v0.4.24
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

WORKDIR /app

COPY services/StateService/go.mod services/StateService/go.sum ./services/StateService/
COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/StateService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /state-service-app ./services/StateService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
COPY --from=builder /state-service-app .
CMD ["./state-service-app"]