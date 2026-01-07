# ---------- Builder ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.work go.sum ./
COPY ./backend/ ./backend/
COPY ./consumer/ ./consumer/
COPY ./libs/ ./libs/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./backend/main.go


# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/app /app/app

USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
