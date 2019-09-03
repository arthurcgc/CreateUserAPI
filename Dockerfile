# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Arthur Coelho <arthur.cavalcante.puc@gmail.com>"

RUN mkdir /app 

ADD . /app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

WORKDIR /app

# Build the Go app
RUN go build -o bin/main src/main.go

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["/app/bin/main"]