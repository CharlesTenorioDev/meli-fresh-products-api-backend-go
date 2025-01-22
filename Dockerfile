FROM golang:latest AS base

WORKDIR /meli-fresh-products-api-backend-t1
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.4

COPY go.mod go.sum ./
RUN go mod download

COPY . .
