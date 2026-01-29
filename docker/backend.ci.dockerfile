FROM gcr.io/distroless/static-debian12:nonroot
EXPOSE 8080

WORKDIR /build
COPY /build/app /build/app
USER nonroot:nonroot

ENTRYPOINT ["/build/app"]
