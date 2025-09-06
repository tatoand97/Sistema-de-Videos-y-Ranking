FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY internal/ ./internal/
COPY cmd/ ./cmd/
COPY assets/ ./assets/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o watermarking ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/watermarking /app/watermarking
COPY --from=builder /app/assets/ /app/assets/
ENV WATERMARK_PATH=/app/assets/nba-logo-removebg-preview.png
CMD ["/app/watermarking"]