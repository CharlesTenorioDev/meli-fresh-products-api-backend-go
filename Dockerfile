FROM golang:1.22.5-alpine as build

WORKDIR /app
COPY . .

RUN go build -o main ./cmd

FROM alpine
WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/db db

EXPOSE 8080
ENTRYPOINT ["./main"]
