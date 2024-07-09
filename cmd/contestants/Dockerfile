FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

# RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
