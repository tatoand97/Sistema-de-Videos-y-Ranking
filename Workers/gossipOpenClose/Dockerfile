FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY gossipOpenClose/go.mod .
COPY shared ./shared
RUN go mod download
COPY gossipOpenClose/internal/ ./internal/
COPY gossipOpenClose/cmd/ ./cmd/
COPY gossipOpenClose/assets/ ./assets/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o gossipopenclose ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/gossipopenclose /app/gossipopenclose
COPY --from=builder /app/assets/ /app/assets/
ENV LOGO_PATH=/app/assets/nba-logo-removebg-preview.png
CMD ["/app/gossipopenclose"]