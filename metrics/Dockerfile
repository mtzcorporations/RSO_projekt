# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /metrics
# Getting go files into image
COPY go.mod .
COPY go.sum .
COPY main.go .
# Installing dependencies
RUN go mod download
# Compiling
RUN go build -o /metrics-compiled

EXPOSE 8005

#CMD [ "go", "run", "main.go" ]
CMD ["/metrics-compiled"]