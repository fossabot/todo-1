FROM golang:alpine

ADD . /go/src/github.com/fharding1/todo
WORKDIR /go/src/github.com/fharding1/todo

RUN apk add --no-cache git

RUN go get -t -v ./...

ENTRYPOINT [ "go", "test", "./..." ]
