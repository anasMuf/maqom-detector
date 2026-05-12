FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api main.go

# ─── Runtime ──────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /bin/api .

EXPOSE 8080

CMD ["./api"]
