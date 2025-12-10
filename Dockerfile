FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o md-serve ./cmd/md-serve

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/md-serve .

EXPOSE 8080

CMD ["./md-serve"]