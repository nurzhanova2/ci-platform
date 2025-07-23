FROM golang:1.24.3-alpine

RUN apk update && apk add --no-cache git docker-cli bash

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o ci-platform ./cmd

EXPOSE 8080

CMD ["./ci-platform"]
