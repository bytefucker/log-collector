
all: clean  build

clean:
	go clean -i ./...
	rm -rf ${GOPATH}/bin/log-collector

build:
	go build -mod=vendor -v -o ${GOPATH}/bin/log-collector .

