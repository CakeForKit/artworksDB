FROM golang:1.24.2-alpine

# (git нужен для go install)
RUN apk add --no-cache git && \
    apk add bash

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN mkdir /migrations
WORKDIR /migrations

RUN pwd && ls -la
COPY . .
RUN ls -la

ENTRYPOINT ["./entrypoint.sh"]
RUN chmod +x ./entrypoint.sh
