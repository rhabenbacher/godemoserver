FROM golang:1.16-alpine

WORKDIR /go/src/training/server
COPY go.mod .
COPY ./server/* ./

RUN GOOS=linux go build -o app/goserver .
CMD ["./app/goserver","standalone"]