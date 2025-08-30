# syntax=docker/dockerfile:1.6

# ---------- Build stage (amd64) ----------
FROM --platform=linux/amd64 golang:1.25-alpine AS build
WORKDIR /src

RUN apk add --no-cache build-base file ca-certificates
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# imprime el entorno de compilación
RUN echo ">>> go env:" && go env GOOS GOARCH CGO_ENABLED

# compila binario 100% estático para linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w -extldflags "-static"' \
    -o /out/main_videork ./cmd/api

# verifica el binario
RUN file /out/main_videork && ls -lh /out/main_videork && sha256sum /out/main_videork

# ---------- Run stage (amd64, Alpine) ----------
FROM --platform=linux/amd64 alpine:3.19 AS run
WORKDIR /app
RUN apk add --no-cache file ca-certificates && update-ca-certificates
COPY --from=build /out/main_videork /app/main_videork
COPY internal/infrastructure/migrations /app/internal/infrastructure/migrations
COPY static /app/static
RUN chmod +x /app/main_videork && file /app/main_videork
EXPOSE 8080
ENTRYPOINT ["/app/main_videork"]
