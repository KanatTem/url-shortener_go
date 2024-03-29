FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /url-shortener ./cmd/url-shortener

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /url-shortener .
COPY config/local.yaml ./config/
COPY migrations ./migrations/

EXPOSE 8082
CMD ["./url-shortener"]