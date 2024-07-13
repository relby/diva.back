FROM golang:1.22.3-alpine AS builder

RUN apk update && \
    apk add --no-cache postgresql-client

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build ./cmd/diva-backups-cron

FROM alpine:latest

WORKDIR /app/

COPY .env ./
COPY --from=builder /app/diva-backups-cron ./

ENTRYPOINT ["./diva-backups-cron"]