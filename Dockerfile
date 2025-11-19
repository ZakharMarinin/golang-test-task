FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ENV CGO_ENABLED=0
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o testovoe ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY migrations ./migrations

COPY --from=builder /app/testovoe .

COPY config.yaml .
COPY .env .env
EXPOSE 8081:8081
CMD ["sh", "-c", "goose -dir ./migrations up postgres \"$GOOSE_DBSTRING\" && ./testovoe"]
