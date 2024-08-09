# syntax=docker/dockerfile:1

FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go .

RUN go build -o /app/bin/server cmd/main.go

RUN go build -o /app/bin/data cmd/data/main.go

RUN go build -o /app/bin/migrate cmd/data/main.go

FROM debian:bullseye-slim

WORKDIR /app

COPY --from=builder /app/bin/server /app/bin/data /app/bin/migrate /app/bin/

EXPOSE 8080

CMD ["/app/bin/server"]
