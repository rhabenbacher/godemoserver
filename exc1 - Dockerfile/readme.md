# Build a go container

## use official golang image 
```
FROM golang:1.16-alpine
```

## create a workdir and copy go files
```
WORKDIR /go/src
COPY ./server/* ./
```
## build goserver
```
RUN go build -o app/goserver .
```

## start goserver as standalone server
```
CMD ["./app/goserver","standalone"]
```
