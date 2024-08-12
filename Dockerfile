# syntax=docker/dockerfile:1

FROM golang:1.22 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/bin/server cmd/main.go

RUN go build -o /app/bin/data cmd/data/*.go

RUN go build -o /app/bin/migrate cmd/migrate/main.go

FROM ubuntu:latest

# install cron and curl
RUN apt-get update && apt-get install -y curl

COPY --from=build /app/bin/server /app/bin/data /app/bin/migrate /app/bin/


RUN apt-get update && apt-get install -y ca-certificates

RUN update-ca-certificates

EXPOSE 8080

CMD ["/app/bin/server"]
