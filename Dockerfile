FROM golang:1.21-alpine as builder
ENV GOPATH /golang

WORKDIR /build
ADD go.mod ./
ADD go.sum ./
ADD cmd ./cmd/
ADD src ./src/
ADD docs ./docs/
RUN go mod download
RUN go build -o ./app-release ./cmd/go-app

# FROM gcr.io/distroless/base-debian11
FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/app-release /app
# COPY public/templates /app/public/templates/
# COPY public/migrations /app/public/migrations/
COPY conf/config.yml /app/conf/config.yml
ENTRYPOINT ["/app/app-release"]
