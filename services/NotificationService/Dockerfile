FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/NotificationService/go.mod ./services/NotificationService/
COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/NotificationService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /notification_service-app ./services/NotificationService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /notification_service-app .
CMD ["./notification_service-app"]