FROM golang:1.15.6-buster AS build-env
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go build -o wordpress_exporter

FROM alpine:latest
WORKDIR /app
RUN apk update && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*
COPY --from=build-env /build/wordpress_exporter /app
ENTRYPOINT ["./wordpress_exporter"]