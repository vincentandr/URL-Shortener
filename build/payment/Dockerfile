# Run using docker compose from root project directory
# Specify environtment file and expose port in docker-compose

FROM golang:1.17-alpine3.15 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY ./cmd/payment/ ./cmd/payment
COPY ./internal/ ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o build ./cmd/payment

# Use scratch for production / only executing binary file
FROM scratch

WORKDIR /app

COPY --from=builder /app/build ./

# If using cmd with executable file as its parameter, means cmd becomes executor of the file (no need entrypoint).
# If using cmd without executable file as its parameter, must use entrypoint for the executor with executable file as entrypoint's parameter
CMD ["./build"]
