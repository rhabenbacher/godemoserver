FROM golang:1.16-alpine as builder

WORKDIR /go/src/training/server
COPY go.mod .
COPY ./server/* /go/src/training/server/


RUN ls -la /go/src/training/*
RUN GOOS=linux go build -o app/goserver .

FROM alpine:latest
WORKDIR /usr/local/bin
COPY --from=builder /go/src/training/server/app/goserver .
CMD ["./goserver","standalone"]