# ---------- Builder ----------
FROM golang:1.25.6-alpine AS builder

WORKDIR /app
RUN apk add --no-cache ca-certificates # git

COPY go.mod go.sum ./
COPY ./transcoder ./transcoder/
COPY ./libs/core/ ./libs/core/
COPY ./libs/storage/ ./libs/storage/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o transcoder ./transcoder/main.go

FROM jrottenberg/ffmpeg:8-alpine
EXPOSE 8080
WORKDIR /app
COPY --from=builder /app/transcoder /app/transcoder
RUN mkdir download output
USER nonroot:nonroot
ENTRYPOINT ["/app/transcoder"]
