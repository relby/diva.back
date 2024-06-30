FROM golang:1.22.3-alpine AS builder

WORKDIR /app/

COPY go.* ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o diva ./cmd/diva/main.go

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/diva/ ./
COPY .env ./

ENV GRPC_HOST=0.0.0.0 GRPC_PORT=50051
EXPOSE 50051

ENTRYPOINT ["./diva"]