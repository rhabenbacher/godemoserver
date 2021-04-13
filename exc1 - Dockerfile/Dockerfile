FROM golang:1.16-alpine

WORKDIR /go/src
COPY * .

RUN go build -o app/goserver .
CMD ["./app/goserver","standalone"]