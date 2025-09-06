FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY StatesMachine/go.mod StatesMachine/go.sum ./
COPY shared ./shared
RUN go mod download
COPY StatesMachine/ .
RUN go mod tidy
RUN go build -o statesmachine ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates netcat-openbsd
WORKDIR /root/
COPY --from=builder /app/statesmachine .
COPY StatesMachine/wait-for-services.sh .
RUN chmod +x wait-for-services.sh
CMD ["./wait-for-services.sh", "./statesmachine"]