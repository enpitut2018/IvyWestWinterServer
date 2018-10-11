FROM golang:1.10.2-alpine3.7

RUN apk add --no-cache git && \
    go get github.com/jinzhu/gorm && \
    go get github.com/gorilla/mux && \
    go get github.com/lib/pq && \
    go get bitbucket.org/liamstask/goose/cmd/goose

WORKDIR /go/src/api
