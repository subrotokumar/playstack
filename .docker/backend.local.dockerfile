FROM golang:1.25.5-alpine

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
COPY cmd/ .
COPY internal/ .
COPY .air.toml .

RUN go mod download

CMD ["air"]