# syntax=docker/dockerfile:1.6

# ---------- Build stage (amd64) ----------
FROM --platform=linux/amd64 golang:1.25-alpine AS build
WORKDIR /src

RUN apk add --no-cache build-base ca-certificates file
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# imprime el entorno de compilación
RUN echo ">>> go env:" && go env GOOS GOARCH CGO_ENABLED

# compila binario 100% estático para linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags='-s -w -extldflags "-static"' \
    -o /out/api ./cmd/api

# verifica el binario
RUN file /out/api && ls -lh /out/api && sha256sum /out/api

# ---------- Run stage (amd64, Alpine) ----------
FROM --platform=linux/amd64 alpine:3.19 AS run
WORKDIR /app
RUN apk add --no-cache file ca-certificates && update-ca-certificates
COPY --from=build /out/api /app/api
COPY internal/infrastructure/migrations /app/internal/infrastructure/migrations

ENV DATABASE_URL=postgres://app_user:app_password@videorank.cfkc2y4acj9g.us-east-1.rds.amazonaws.com:5432/videorank?sslmode=require \
    JWT_SECRET=6EPOc2/6veZ71FAJSF68iv21ho83NLQSaycHEGTTjGO9TBmRsphMp5JqgieFTcGn \
    PORT=8080 \
    CORS_ORIGIN=* \
    UPLOAD_DIR=/uploads \
    REDIS_ADDR=10.0.3.71:6379 \
    CACHE_PREFIX=videorank: \
    CACHE_TTL_SECONDS=120 \
    AWS_REGION=us-east-1 \
    AWS_ACCESS_KEY_ID=ASIAXGQWHVU2B6IPJZIB \
    AWS_SECRET_ACCESS_KEY=CRXbBkoAJus5XyYnBRJO4uh3NDG7BCawjB95a73a \
    S3_BUCKET=video-ranking-raw-videos \
    S3_ENDPOINT=\
    S3_USE_PATH_STYLE=false \
    RABBITMQ_URL=amqp://admin:admin@10.0.5.20:5672/ \
    STATES_MACHINE_QUEUE=states_machine_queue \
    RABBITMQ_QUEUE_MAXLEN=1000

RUN chmod +x /app/api && file /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]

