FROM golang:1.24.1-alpine

RUN mkdir /build
WORKDIR /build

COPY ../../go.* .
RUN go mod download

CMD ["make", "run-app"]
