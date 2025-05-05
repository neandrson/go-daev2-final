# syntax=docker/dockerfile:1

FROM golang:1.23.1-alpine AS builder

# Set destination for COPY
WORKDIR /

COPY . .

# Download Go modules
COPY go.* ./
RUN go mod download

# Build
RUN GOOS=linux go build -o agent ./cmd/main.go

FROM alpine:3.18

WORKDIR /

COPY --from=builder . .

CMD ["./agent"]