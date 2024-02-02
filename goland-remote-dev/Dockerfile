FROM golang:1.21.5-bullseye
ENV GOFLAGS=-buildvcs=false
ARG USERNAME=go
ARG GROUPNAME=go
ARG UID=1000
ARG GID=1000
ARG IDEURL=https://download.jetbrains.com/go/goland-2023.3.3.tar.gz
ENV TZ Asia/Tokyo
ENV APP_ROOT /home/$USERNAME/app
RUN apt-get update && \
    apt-get install -y vim && \
    rm -rf /var/lib/apt/lists/*
RUN groupadd -g $GID $GROUPNAME && \
   useradd -m -s /bin/bash -u $UID -g $GID $USERNAME

WORKDIR $APP_ROOT
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download
COPY goland-remote-dev .
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
ENV EDITOR vim
RUN curl -fsSL -o ide.tar.gz $IDEURL && \
rm -rf ide && \
mkdir ide && \
tar xfz ide.tar.gz --strip-components=1 -C ide && \
rm ide.tar.gz
CMD ["sh", "-c", "air -c .air.toml & ide/bin/remote-dev-server.sh run $APP_ROOT --listenOn 0.0.0.0 --port 5994"]
EXPOSE 8080
EXPOSE 5994