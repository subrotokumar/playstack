FROM golang:1.25.5-alpine

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ .
COPY internal/ .
COPY .air.toml .
CMD ["air"]