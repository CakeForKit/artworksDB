
FROM golang:1.24.2-alpine

RUN mkdir /build
WORKDIR /build

COPY ../../go.* .
RUN go mod download

COPY ../../. .

CMD ["go", "run", "./cmd/app/main.go"]

# CMD ["go", "test"]
