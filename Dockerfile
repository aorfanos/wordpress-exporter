FROM golang:1.15.6-buster AS build-env
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go build -o wordpress_exporter

FROM debian:buster-slim
WORKDIR /app
COPY --from=build-env /build/wordpress_exporter /app
ENTRYPOINT ["./wordpress_exporter"]