FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY Watermarking/go.mod .
COPY shared ./shared
RUN sed -i 's|replace shared => ../shared|replace shared => ./shared|' go.mod
RUN go mod download
COPY Watermarking/internal/ ./internal/
COPY Watermarking/cmd/ ./cmd/
COPY Watermarking/assets/ ./assets/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o watermarking ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/watermarking /app/watermarking
COPY --from=builder /app/assets/ /app/assets/
ENV WATERMARK_PATH=/app/assets/nba-logo-removebg-preview.png
CMD ["/app/watermarking"]