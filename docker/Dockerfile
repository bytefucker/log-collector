
#build stage
FROM golang:1.13.5-alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ./app
LABEL Name=docker Version=0.0.1

