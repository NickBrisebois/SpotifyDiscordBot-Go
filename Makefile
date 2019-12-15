VERSION=$(shell git describe --tags)
BUILD=$(shell git rev-parse --short HEAD)

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=govendor fetch
VENDORINIT=govendor init
BINARY_NAME=spottybot
VERSION=$(shell cat ./VERSION)

LDFLAGS=-ldflags "-X=main.Version=$(VESRION) -X=main.Build=$(BUILD)"

all: build

.PHONY: build
build:
	rm -rf ./build/;
	mkdir ./build;
	cp -r ./config/config.toml ./build/
	$(MAKE) -s go-build

test:
	$(GTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf ./build/

go-build:
	@GOPATH=$(GOPATH) go build $(LDFLAGS) -o ./build/$(BINARY_NAME)

run:
	./build/$(BINARY_NAME) --config ./build/config.toml

build-docker:
	docker build -t $(BINARY_NAME)-$(VERSION) .

deps:
	go get -u github.com/kardianos/govendor
	$(VENDORINIT)
	$(GOGET) github.com/BurntSushi/toml
	$(GOGET) github.com/bwmarrin/discordgo
	$(GOGET) github.com/zmb3/spotify
	$(GOGET) golang.org/x/oauth2/clientcredentials
	$(GOGET) github.com/mvdan/xurls
	$(GOGET) go.mongodb.org/mongo-driver/mongo
	$(GOGET) go.mongodb.org/mongo-driver/bson
	$(GOGET) go.mongodb.org/mongo-driver/mongo/options

