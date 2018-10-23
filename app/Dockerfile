FROM golang:1.10.2-alpine3.7

RUN apk add --no-cache git gcc && \
  go get github.com/jinzhu/gorm && \
  go get github.com/gorilla/mux && \
  go get github.com/lib/pq &&\
  go get github.com/rs/xid && \
  go get github.com/aws/aws-sdk-go

WORKDIR /go/src/github.com/enpitut2018/IvyWestWinterServer
