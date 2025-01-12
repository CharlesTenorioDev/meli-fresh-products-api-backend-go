FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . .

RUN go install github.com/air-verse/air@latest
RUN go mod download

CMD ["air", "-c", ".air.toml"]
