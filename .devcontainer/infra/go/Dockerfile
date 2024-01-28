FROM golang:1.20.0-buster

ENV TZ=Asia/Tokyo
WORKDIR /data

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest