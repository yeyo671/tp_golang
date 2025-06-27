# Build stage
FROM golang:1.24.4-alpine AS builder

# Installer les dépendances pour CGO
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Construire avec CGO activé
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o urlshortener .

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/urlshortener .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./urlshortener", "run-server"]