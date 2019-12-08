GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goproj

all: build
build:
	rm -rf ./build/
	mkdir ./build
	cp -r ./src/config/* ./build/
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v ./cmd/ ./internal/ 

test:
	$(GTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf ./build/

run:
	./build/$(BINARY_NAME)

deps:
	$(GOGET) github.com/BurnSushi/toml

