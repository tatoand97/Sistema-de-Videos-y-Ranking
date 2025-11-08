FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY TrimVideo/go.mod .
COPY shared ./shared
RUN sed -i 's|replace shared => ../shared|replace shared => ./shared|' go.mod
RUN go mod download
COPY TrimVideo/internal/ ./internal/
COPY TrimVideo/cmd/ ./cmd/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o trimvideo ./cmd/main.go

FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/trimvideo /app/trimvideo
CMD ["/app/trimvideo"]