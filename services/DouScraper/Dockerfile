FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY services/DouScraper/go.mod ./services/DouScraper/

COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/services/DouScraper
RUN go mod download

WORKDIR /app
COPY . .

RUN go build -o /dou_scraper-app ./services/DouScraper/cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /dou_scraper-app .
CMD ["./dou_scraper-app"]