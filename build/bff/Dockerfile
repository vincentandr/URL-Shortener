# Run using docker compose from root project directory
# Specify environtment file and expose port in docker-compose

FROM golang:1.17-alpine3.15 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

# Cache the mod file, only redownload if it changes
RUN go mod download

COPY ./cmd/bff ./cmd/bff
COPY ./internal ./internal

# CGO should be enabled if the codes use C library. GOOS is image OS
RUN CGO_ENABLED=0 GOOS=linux go build -o build ./cmd/bff

# Use scratch for production / only executing binary file
FROM scratch

WORKDIR /app

# Copy resulted built app from previous stage
COPY --from=builder /app/build ./

# If using cmd with executable file as its parameter, means cmd becomes executor of the file (no need entrypoint).
# If using cmd without executable file as its parameter, must use entrypoint for the executor with executable file as entrypoint's parameter
CMD ["./build"]

