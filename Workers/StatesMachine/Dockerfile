FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY StatesMachine/go.mod StatesMachine/go.sum ./
COPY shared ./shared
COPY StatesMachine/ .
RUN sed -i 's|replace shared => ../shared|replace shared => ./shared|' go.mod
RUN go mod download
RUN go mod tidy
RUN go build -o statesmachine ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates netcat-openbsd
WORKDIR /root/
COPY --from=builder /app/statesmachine .
COPY --from=builder /app/wait-for-services.sh .

CMD ["/app/statesmachine"]