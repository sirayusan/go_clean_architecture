FROM golang:1.20.0-buster

# Task Runner のインストール
RUN curl -sL "https://github.com/go-task/task/releases/download/v3.9.0/task_linux_amd64.deb" -o task.deb && dpkg -i task.deb && rm task.deb

ENV TZ=Asia/Tokyo
WORKDIR /data

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/codegangsta/gin@latest
RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/fatih/gomodifytags@latest