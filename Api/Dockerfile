# Build stage
FROM golang:1.24-alpine AS build
WORKDIR /Api
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main_videork ./cmd/api

# Run stage
FROM alpine:3.19

COPY --from=build /Api/main_videork ./main_videork
COPY internal/infrastructure/migrations internal/infrastructure/migrations
COPY static static
EXPOSE 8080
CMD ["./main_videork"]