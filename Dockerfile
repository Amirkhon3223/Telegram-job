FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Build Bot
RUN CGO_ENABLED=0 GOOS=linux go build -o /bot ./cmd/bot

# API image
FROM alpine:3.19 AS api
RUN apk --no-cache add ca-certificates
COPY --from=builder /api /api
CMD ["/api"]

# Bot image
FROM alpine:3.19 AS bot
RUN apk --no-cache add ca-certificates
COPY --from=builder /bot /bot
CMD ["/bot"]
