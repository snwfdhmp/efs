FROM golang:1.13.7-alpine3.11 as builder
ENV GO111MODULE=on

RUN apk update && apk add git

WORKDIR /go/src/github.com/snwfdhmp/efs

COPY ./server ./server


RUN go get ./...
RUN go build -o ./server ./server/cmd/server

FROM alpine:3.11

WORKDIR /server

COPY --from=builder /go/src/github.com/snwfdhmp/efs/server efs-server

RUN chmod go-rwx ./efs-server

ENTRYPOINT [ "/server/efs-server" ]