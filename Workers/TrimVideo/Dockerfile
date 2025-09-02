# ---------- builder ----------
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

# Descarga lo que ya esté en go.mod (cache-friendly)
COPY go.mod .
RUN go mod download

# Copia el código y resuelve deps reales
COPY internal/ ./internal/
COPY cmd/ ./cmd/
RUN go mod tidy

# Compila binario estático
RUN CGO_ENABLED=0 go build -o trimvideo ./cmd/main.go

# ---------- runtime ----------
FROM alpine:3.19
RUN apk add --no-cache ffmpeg ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /app/trimvideo /app/trimvideo
CMD ["/app/trimvideo"]