FROM golang:alpine

WORKDIR /home/container/ypn

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY / /home/container/ypn

RUN export GIN_MODE=release

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go build -o ypn

EXPOSE 9000

ENTRYPOINT ["./ypn"]