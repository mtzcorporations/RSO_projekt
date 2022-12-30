# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /health
# Getting go files into image
COPY go.mod .
COPY go.sum .
COPY main.go .
# Installing dependencies
RUN go mod download
# Compiling
RUN go build -o /health-compiled

EXPOSE 8080

CMD ["/health-compiled"]