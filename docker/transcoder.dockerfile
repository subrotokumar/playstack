# ---------- Builder ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.mod go.sum ./
COPY ./transcoder ./transcoder/
COPY ./libs/core/ ./libs/core/
COPY ./libs/storage/ ./libs/storage/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./transcoder/main.go

FROM jrottenberg/ffmpeg:8-alpine

WORKDIR /app
COPY --from=builder /app/app /app/transcoder
RUN mkdir download output
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]
