FROM golang:1.25.1-alpine AS builder

WORKDIR /app
COPY ./go.* .
RUN go mod download
# Копируем весь проект (из корня, так как context: ../)
COPY . ./       
RUN go build -o /app/my-app ./cmd/app


FROM golang:1.25.1-alpine AS test-runner

# Для healthcheck или скриптов
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY . .
# Копируем скомпилированное приложение из стадии builder (для E2E тестов) --ПОТОМ
# COPY --from=builder /app/my-app ./my-app
RUN go mod download
CMD ["tail", "-f", "/dev/null"]

