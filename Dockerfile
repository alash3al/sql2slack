FROM golang:alpine

# RUN apk add --no-cache git

RUN CGO_ENABLED=0 go get -ldflags="-s -w" github.com/alash3al/sql2slack

ENTRYPOINT ["sql2slack"]

WORKDIR /root/
