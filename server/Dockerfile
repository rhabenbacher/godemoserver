FROM golang:1.16-alpine as builder

WORKDIR /go/src
COPY * .
RUN go build -o app/goserver .

FROM alpine:latest
WORKDIR /usr/local/bin
COPY --from=builder /go/src/app/goserver .
CMD ["./goserver","standalone"]
