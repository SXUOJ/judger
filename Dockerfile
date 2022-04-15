FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive \
    GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct\ 
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# apt
RUN apt-get update && \
    apt-get -y install gcc golang-go && \ 
    apt-get -y install git cmake libseccomp-dev 

# judger
RUN git clone https://github.com/SXUOJ/judger.git /tmp/judger && \
    cd /tmp/judger && go build -o /judger 

# clear
RUN    rm -rf /tmp/judger && \ 
    apt-get purge -y --auto-remove cmake git && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

EXPOSE 9000

ENTRYPOINT ["/judger"]