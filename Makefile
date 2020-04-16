
all: clean  build

clean:
	go clean -i .

build:
	go build -mod=vendor -v -o ${GOPATH}/bin/logcollect/agent ./agent/cmd/agent
	go build -mod=vendor -v -o ${GOPATH}/bin/logcollect/analysis ./analysis/cmd/analysis
	go build -mod=vendor -v -o ${GOPATH}/bin/logcollect/manager ./manager/cmd/manager

