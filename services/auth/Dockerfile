FROM golang:1.18.0-alpine3.15 AS build
WORKDIR /go/src/github.com/auth
COPY . .
RUN go build -o auth
CMD ["./auth"]
EXPOSE 3004
EXPOSE 4004