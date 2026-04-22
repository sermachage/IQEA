# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /build

COPY . .

RUN go mod download && \
    go build -o demographic-api .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /build/demographic-api .

EXPOSE 8080

ENV PORT=8080

CMD ["./demographic-api"]
