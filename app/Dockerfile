FROM golang:1.10.2-alpine3.7

RUN apk add --no-cache git gcc && \
  go get github.com/jinzhu/gorm && \
  go get github.com/gorilla/mux && \
  go get github.com/lib/pq

WORKDIR /go/src/github.com/enpitut2018/IvyWestWinterServer
