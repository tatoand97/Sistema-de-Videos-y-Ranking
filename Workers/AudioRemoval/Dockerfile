FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY AudioRemoval/go.mod AudioRemoval/go.sum ./
COPY shared ./shared
COPY AudioRemoval/internal/ ./internal/
COPY AudioRemoval/cmd/ ./cmd/
RUN sed -i 's|replace shared => ../shared|replace shared => ./shared|' go.mod
RUN go mod download
RUN go mod tidy
RUN go build -o audioremoval cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates ffmpeg
RUN mkdir -p /tmp && chmod 777 /tmp

WORKDIR /app
COPY --from=builder /app/audioremoval .

USER nobody
CMD ["./audioremoval"]