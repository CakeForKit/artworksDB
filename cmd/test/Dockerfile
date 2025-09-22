# FROM golang:1.25.1-alpine AS test-runner

# RUN apk add --no-cache git bash curl 
# WORKDIR /app

# COPY ./go.* .
# RUN go mod download
# COPY . .

# RUN go build -o ./test-app ./cmd/test/main.go
# CMD ["./test-app"]

FROM golang:1.25.1-alpine AS test-runner

RUN apk add git bash 
RUN go install github.com/air-verse/air@latest
WORKDIR /app

COPY ./go.* .
RUN go mod download
COPY . .

RUN mkdir -p tmp

CMD ["air", "-c", "./cmd/test/.air.toml"]

# CMD ["tail", "-f", "/dev/null"]

