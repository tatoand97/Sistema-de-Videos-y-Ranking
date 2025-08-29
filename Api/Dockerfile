# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-api ./cmd/api

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/auth-api ./auth-api
COPY --from=build /app/internal/infrastructure/migrations ./internal/infrastructure/migrations
EXPOSE 8080
CMD ["./auth-api"]