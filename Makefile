# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


all: clean test build

clean:
	$(GOCLEAN)

test:
	$(GOTEST) -v ./...

build:
	$(GOBUILD)  -v ./...

