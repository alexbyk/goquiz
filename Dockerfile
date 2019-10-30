FROM golang:1.13-alpine as builder

LABEL maintainer="alexbyk <alex@alexbyk.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o integrator ./cmd/integrator/main.go
RUN go build -o parser ./cmd/parser/main.go
RUN go build -o httpapi ./cmd/httpapi/main.go

######## Start a new stage from scratch #######
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/integrator /app/parser /app/httpapi /app/wait-for /app/data.csv ./
