FROM golang:1.13
MAINTAINER Oliver Nadj <mr.oliver.nadj@gmail.com>
RUN go get -u github.com/gorilla/mux
ADD . /go/src/goapi
ADD ./go /go/src/github.com/olivernadj/post-proc/goapi/go
EXPOSE 8080
WORKDIR /go/src/goapi
ENTRYPOINT [ "go", "run",  "main.go"]