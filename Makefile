GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=govendor fetch
BINARY_NAME=spoticord

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
	$(GOGET) init
	$(GOGET) github.com/BurntSushi/toml
	$(GOGET) github.com/bwmarrin/discordgo

