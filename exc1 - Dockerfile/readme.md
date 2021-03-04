# Build a go container

## use official golang image 
```
FROM golang:1.16-alpine
```

## create a workdir and copy go files
```
WORKDIR /go/src/training/server
COPY go.mod .
COPY ./server/* /go/src/training/server/
```
## build goserver
```
RUN GOOS=linux go build -o app/goserver .
```

## start goserver as standalone server
```
CMD ["./app/goserver","standalone"]
```
