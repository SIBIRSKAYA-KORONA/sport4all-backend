#builder
FROM golang:alpine3.11

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /application

COPY . .

RUN make build-server