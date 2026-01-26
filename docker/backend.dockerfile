# ---------- Builder ----------
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.mod go.sum ./
COPY ./backend/ ./backend/
COPY ./libs/core/ ./libs/core/
COPY ./libs/idp/ ./libs/idp/
COPY ./libs/db/ ./libs/db/
COPY ./libs/storage/ ./libs/storage/

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o ./tmp/backend ./backend/main.go

# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /app/tmp/backend /app/backend

USER nonroot:nonroot
ENTRYPOINT ["/app/backend"]
