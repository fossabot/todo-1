FROM golang:latest

ADD . /go/src/github.com/fharding1/todo
WORKDIR /go/src/github.com/fharding1/todo

RUN go get -t -v ./...

ENTRYPOINT [ "go", "test", "./..." ]
