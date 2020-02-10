FROM golang:1.13.7-alpine3.11 as builder
ENV GO111MODULE=on

RUN apk update && apk add git

WORKDIR /go/src/github.com/snwfdhmp/efs/server

COPY ./ ./

RUN go get ./...

RUN go build -o ./dist-server ./cmd/server

FROM alpine:3.11

WORKDIR /server

COPY --from=builder /go/src/github.com/snwfdhmp/efs/server/dist-server efs-server

RUN chmod go-rwx ./efs-server

ENTRYPOINT [ "/server/efs-server" ]