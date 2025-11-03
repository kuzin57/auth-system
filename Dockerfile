FROM golang:1.24.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o main cmd/main/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssl

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/init.sh ./init.sh

RUN chmod +x ./init.sh

RUN mkdir -p config

ENTRYPOINT ["./init.sh"]
CMD ["./main", "-config", "config/config.yaml", "-secrets", "config/secrets.yaml"]
