FROM golang:alpine
LABEL authors="ricar"

WORKDIR /app
COPY . .
RUN go mod tidy && go build -o api .

EXPOSE 8443
CMD ["./api"]