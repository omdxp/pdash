FROM golang:1.18.0-alpine3.15 AS build
WORKDIR /go/src/github.com/orders
COPY . .
RUN go build -o orders
CMD ["./orders"]
EXPOSE 3002
EXPOSE 4002