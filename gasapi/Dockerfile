FROM golang:1.18-alpine

WORKDIR /gasapi
# Getting go files into image
COPY go.mod .
COPY go.sum .
COPY main.go .
# Installing dependencies
RUN go mod download
# Compiling
RUN go build -o /gasapi-compiled

EXPOSE 8004

CMD ["/gasapi-compiled"]
