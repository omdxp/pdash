FROM golang:1.18.0-alpine3.15 AS build
WORKDIR /go/src/github.com/suppliers
COPY . .
RUN go build -o suppliers
CMD ["./suppliers"]
EXPOSE 3003
EXPOSE 4003