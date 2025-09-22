FROM golang:1.25.1-alpine AS test-runner

RUN apk add --no-cache git bash curl 
# postgresql-client
WORKDIR /app

COPY ./go.* .
RUN go mod download
COPY . .

RUN go build -o ./test-app ./cmd/test/main.go

CMD ["./test-app"]

# FROM golang:1.25.1-alpine AS test-runner

# RUN apk add git bash 
# RUN go install github.com/air-verse/air@latest
# WORKDIR /app

# COPY ./go.* .
# RUN go mod download
# COPY . .

# RUN mkdir -p tmp
# RUN ls -la
# COPY ./cmd/app/.air.toml .
# RUN chmod 644 .air.toml
# RUN ls -la ./cmd/test
# RUN ls -la

# RUN cat ./cmd/app/.air.toml
# RUN cat .air.toml
# # CMD ["air", "-c", ".air.toml"]

# # CMD ["tail", "-f", "/dev/null"]

