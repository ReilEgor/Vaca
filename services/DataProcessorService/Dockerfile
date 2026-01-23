FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/DataProcessorService/go.mod ./services/DataProcessorService/

COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/DataProcessorService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /data_processor_service-app ./services/DataProcessorService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /data_processor_service-app .
CMD ["./data_processor_service-app"]