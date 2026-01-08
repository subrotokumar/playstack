FROM golang:1.25.5-alpine

WORKDIR /app

RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest && \
    go install github.com/go-task/task/v3/cmd/task@latest

COPY go.mod go.sum ./
COPY ./backend/ ./backend/
COPY ./libs/core/ ./libs/core/
COPY ./libs/core/ ./libs/core/
COPY ./libs/idp/ ./libs/idp/
COPY .air.toml .
CMD ["air"]