

# --- Etapa de Compilación Go ---
FROM golang:1.21-alpine AS go-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY cmd/ ./cmd/
RUN go mod tidy
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o curtaininjector cmd/main.go

# --- Etapa Final ---
FROM alpine:latest

RUN apk add --no-cache ca-certificates ffmpeg
RUN mkdir -p /tmp && chmod 777 /tmp

WORKDIR /app

COPY ANB_Curtain_IN.mp4 /app/ANB_Curtain_IN.mp4
COPY ANB_Curtain_OUT.mp4 /app/ANB_Curtain_OUT.mp4
COPY --from=go-builder /app/curtaininjector .

USER nobody
CMD ["./curtaininjector"]