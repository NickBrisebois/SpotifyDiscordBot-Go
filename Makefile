GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=govendor fetch
VENDORINIT=govendor init
BINARY_NAME=spottybot

all: build
build:
	rm -rf ./build/;
	mkdir ./build;
	cp -r ./config/config.toml ./build/
	$(GOBUILD) -o ./build/$(BINARY_NAME) -v

test:
	$(GTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf ./build/

run:
	./build/$(BINARY_NAME) --config ./build/config.toml

deps:
	go get -u github.com/kardianos/govendor
	$(VENDORINIT)
	$(GOGET) github.com/BurntSushi/toml
	$(GOGET) github.com/bwmarrin/discordgo
	$(GOGET) github.com/zmb3/spotify
	$(GOGET) golang.org/x/oauth2/clientcredentials

