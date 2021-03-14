FROM golang:1.16-alpine as builder

WORKDIR /go/src
COPY ./server/* ./

RUN GOOS=linux go build -o app/goserver .
RUN find /go/src

FROM alpine:latest
EXPOSE 8080
WORKDIR /usr/local/bin
COPY --from=builder /go/src/app/goserver .
CMD ["./goserver","standalone","--port=8080"]
