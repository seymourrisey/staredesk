FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o staredesk ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/staredesk .

EXPOSE 8080

CMD ["./staredesk"]
