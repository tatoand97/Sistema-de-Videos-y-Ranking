# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates

# Capa de módulos (rápida en rebuilds)
COPY go.mod ./
RUN go mod download

# Copia el código y resuelve dependencias indirectas → genera go.sum
COPY . .
RUN go mod tidy

# Compila
RUN CGO_ENABLED=0 go build -o /bin/admincache ./cmd/admincache

# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/admincache /app/admincache
USER nonroot:nonroot
ENTRYPOINT ["/app/admincache"]