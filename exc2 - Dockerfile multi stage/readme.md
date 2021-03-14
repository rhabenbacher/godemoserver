# Build a go container - multi stage

## use official golang image to build
```yml
FROM golang:1.16-alpine as builder
```
##
## create a workdir and copy go files
```yml
WORKDIR /go/src
COPY ./server/* ./
```
## build goserver
```yml
RUN GOOS=linux go build -o app/goserver .
```
## select alpine for production use
```yml
FROM alpine:latest
```
## create a workdir
```yml
WORKDIR /usr/local/bin
```
## copy executable from builder stage
```yml
COPY --from=builder /go/src/app/goserver .
```
## start goserver in standalone mode
```yml
CMD ["./goserver","standalone"]
```