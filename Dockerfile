FROM golang:alpine as builder

RUN apk add --no-cache git

RUN CGO_ENABLED=0 go get -ldflags="-s" github.com/alash3al/sql2slack

FROM alpine

COPY --from=builder /go/bin/sql2slack ./

ENTRYPOINT ["./sql2slack"]
