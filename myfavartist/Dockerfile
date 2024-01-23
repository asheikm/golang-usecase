FROM golang:1.16-alpine as builder
WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
