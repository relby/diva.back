FROM golang:1.22.3-alpine

RUN apk update && apk upgrade && apk add bash && rm -rf /var/cache/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.20.0/goose_linux_x86_64 /usr/bin/goose

RUN chmod +x /usr/bin/goose

WORKDIR /migrations/

COPY ./database/postgres/migrations/ ./

ENTRYPOINT ["/usr/bin/goose", "postgres", "up"]