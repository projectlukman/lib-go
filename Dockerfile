FROM golang:1.16-alpine
LABEL maintainer="project.lukman@gmail.com"

RUN apk update && \
    apk add bash git && \
    apk add gcc && \
    apk add musl-dev && \
    apk add curl && \
    apk add --update make

ENV TZ=Asia/Jakarta
RUN apk add -U tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY . /home/golang/lib
WORKDIR /home/golang/lib