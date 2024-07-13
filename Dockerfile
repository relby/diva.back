FROM golang:1.22.3-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN go build ./cmd/diva/
RUN go build ./cmd/diva-backups-cron/
RUN go build ./cmd/diva-create-admin/
RUN go build ./cmd/diva-excel/

FROM alpine:latest

ENV GRPC_HOST=0.0.0.0 GRPC_PORT=50051

WORKDIR /app/

COPY .env ./
COPY --from=builder /app/diva ./
COPY --from=builder /app/diva-backups-cron ./
COPY --from=builder /app/diva-create-admin ./
COPY --from=builder /app/diva-excel ./

EXPOSE 50051

ENTRYPOINT ["./diva"]