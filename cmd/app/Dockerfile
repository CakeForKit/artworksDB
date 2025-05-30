# Стадия сборки (builder) - это "временный слой"
FROM golang:1.24.2-alpine AS builder

# Устанавливаем зависимости (добавляем swag)
RUN apk add --no-cache git bash && \
    go install github.com/swaggo/swag/cmd/swag@latest

RUN mkdir /build
WORKDIR /build

# Копируем только файлы модулей сначала для эффективного кэширования
COPY ./go.* .
RUN go mod download

# Копируем весь проект (из корня, так как context: ../)
COPY . .

# RUN CGO_ENABLED=0 GOOS=linux go list -m all
# Перемещаем docs в правильное место (чтобы соответствовать пути модуля)
# RUN mkdir -p /build/src/git.iu7.bmstu.ru/ped22u691/PPO.git && \
    # mv docs /build/src/git.iu7.bmstu.ru/ped22u691/PPO.git/

# Собираем приложение с оптимизациями
# CGO_ENABLED=0 - Генерирует полностью статический бинарный файл
# -ldflags="-s -w" - Удаляет отладочную и другую необязательную информацию
RUN /go/bin/swag init -g ./cmd/app/main.go --output ./docs && \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /build/bin/artwork-app \
    ./cmd/app/main.go
# RUN CGO_ENABLED=0 GOOS=linux go build \
#     -ldflags="-s -w" \                  
#     -o /build/bin/artwork-app \
#     ./cmd/app/main.go


# Финальная стадия - Новый базовый слой
FROM alpine:latest  

RUN mkdir /app
WORKDIR /app

# Копируем ТОЛЬКО бинарник из builder  
COPY --from=builder /build/bin/artwork-app .

# Копируем необходимые файлы конфигурации
COPY --from=builder /build/configs ./configs
COPY --from=builder /build/migrations ./migrations
# COPY --from=builder /build/src/git.iu7.bmstu.ru/ped22u691/PPO.git/docs ./docs

# Устанавливаем зависимости (добавляем swag)
RUN apk add --no-cache git bash && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Устанавливаем зависимости для runtime (если нужны)
RUN apk add --no-cache ca-certificates tzdata

# Порт, который слушает приложение
EXPOSE 8080

# Команда запуска (GIN_MODE будет переопределен из docker-compose)
CMD ["./artwork-app"]


# FROM golang:1.24.2-alpine

# RUN mkdir /build
# WORKDIR /build

# COPY ../../go.* .
# RUN go mod download

# COPY ../../. .

# RUN go build -o app ./cmd/app/main.go
# EXPOSE 8080
# CMD ["./app"]

# # CMD ["go", "run", "./cmd/app/main.go"]

# # CMD ["go", "test"]
