FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/CoordinatorService/go.mod ./services/CoordinatorService/

COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/CoordinatorService
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /coordinator_service-app ./services/CoordinatorService/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /coordinator_service-app .
CMD ["./coordinator_service-app"]