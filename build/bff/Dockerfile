FROM golang:1.17-alpine3.15

WORKDIR /app/

COPY go.mod .
COPY go.sum .

# Cache the mod file, only redownload if it changes
RUN go mod download

COPY cmd/bff .
COPY cmd/internal .

# GOOS is image OS
RUN GOOS=linux go build 