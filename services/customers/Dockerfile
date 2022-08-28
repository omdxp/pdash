FROM golang:1.18.0-alpine3.15 AS build
WORKDIR /go/src/github.com/customers
COPY . .
RUN go build -o customers
CMD ["./customers"]
EXPOSE 3001
EXPOSE 4001