
all: clean  build

clean:
	go clean -i ./...
	rm ${GOPATH}/bin/logcollect

build:
	go build -mod=vendor -v -o ${GOPATH}/bin/logcollect ./cmd/logcollect

