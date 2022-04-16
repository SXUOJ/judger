FROM golang:alpine AS builder
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct\ 
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN  go build -o judger .


FROM ubuntu:20.04
ENV DEBIAN_FRONTEND=noninteractive 

# apt
RUN apt-get upgrade && apt-get update && \
    apt-get -y install gcc golang-go && \ 
    apt-get -y install git cmake libseccomp-dev 

# clear
RUN apt-get purge -y --auto-remove cmake git && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/judger /

EXPOSE 9000

ENTRYPOINT ["/judger"]