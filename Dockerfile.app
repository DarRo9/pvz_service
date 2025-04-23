FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM alpine:latest
WORKDIR /root/

RUN mkdir -p ./config

COPY --from=builder /server .
COPY --from=builder /app/config/config.yaml ./config/

ENTRYPOINT [ "./server" ]