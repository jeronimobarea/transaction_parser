FROM golang:1.24.2-alpine AS builder

WORKDIR /src/app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o transaction_parser cmd/main.go

FROM alpine:latest AS release

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /src/app/transaction_parser .

CMD ["./transaction_parser"]
