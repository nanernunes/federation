FROM      golang:1.17-alpine3.15 AS builder

RUN       apk add --update make
RUN       go install github.com/swaggo/swag/cmd/swag@latest

ENV       APP=federation
ENV       CGO_ENABLED=0

ARG       DIR=/${APP}

WORKDIR   $DIR
ADD .     $DIR

RUN       make deps build



FROM      alpine:latest
COPY      --from=builder /federation/federation /usr/local/bin/federation

EXPOSE    80
CMD       federation
