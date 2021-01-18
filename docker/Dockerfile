
#build stage
FROM golang:1.13.5-alpine AS builder
WORKDIR /go/src/logcollect
COPY . .
RUN echo -e "http://mirrors.aliyun.com/alpine/v3.11/main\nhttp://mirrors.aliyun.com/alpine/v3.11/community" > /etc/apk/repositories
RUN apk add --no-cache git
RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go install -mod=vendor -v ./...

#final stage
FROM alpine:3.11
RUN echo -e "http://mirrors.aliyun.com/alpine/v3.11/main\nhttp://mirrors.aliyun.com/alpine/v3.11/community" > /etc/apk/repositories
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/* /log-collector/


