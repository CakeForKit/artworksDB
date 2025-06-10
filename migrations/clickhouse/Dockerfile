FROM golang:1.24.2-alpine

# Устанавливаем зависимости
RUN apk add --no-cache git bash curl 

# Устанавливаем migrate с поддержкой ClickHouse
RUN go install -tags 'clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Устанавливаем clickhouse-client через curl (альтернативный способ)
RUN curl https://clickhouse.com/ | sh && \
    mv ./clickhouse /usr/local/bin/

# Копируем скрипты миграции
RUN mkdir /clickhouse_migrations
WORKDIR /clickhouse_migrations
COPY . .

# Делаем скрипт исполняемым
RUN chmod +x ./entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]