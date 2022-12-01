# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /authentication
# Getting go files into image
COPY go.mod .
COPY go.sum .
COPY main.go .

# Installing dependencies
RUN go mod download
# Compiling
RUN go build -o /authentication-compiled

EXPOSE 8002

#CMD [ "go", "run", "main.go" ]
CMD ["/authentication-compiled"]
