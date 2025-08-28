FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
COPY internal/ ./internal/
COPY cmd/ ./cmd/
RUN go mod tidy
RUN go mod download
RUN go build -o audioremoval cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN mkdir -p /tmp && chmod 777 /tmp

WORKDIR /app
COPY --from=builder /app/audioremoval .
COPY .env .

USER nobody
CMD ["./audioremoval"]