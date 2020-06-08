# Start from the latest golang base image
FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="Madhan Raj <jmadhanraj96@gmail.com>"

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY ./src .

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]
