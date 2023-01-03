# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /posts
# Getting go files into image
COPY go.mod .
COPY go.sum .
COPY main.go .

# Installing dependencies
RUN go mod download
# Compiling
RUN go build -o /posts-compiled

EXPOSE 8003

#CMD [ "go", "run", "main.go" ]
CMD ["/posts-compiled"]
