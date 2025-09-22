
FROM golang:1.25.1-alpine AS test-runner

RUN apk add git bash 
RUN go install github.com/air-verse/air@latest
WORKDIR /app

COPY ./go.* .
RUN go mod download
COPY . .

RUN mkdir -p tmp
# COPY ./cmd/app/.air.toml .
RUN ls -la 
RUN tree
CMD ["air", "-c", "./cmd/app/.air.toml"]


# CMD ["tail", "-f", "/dev/null"]

