FROM golang:1.22.5-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o cmd

FROM alpine
WORKDIR /app

COPY --from=build /app/main .

EXPOSE 8080
ENTRYPOINT ["main"]
