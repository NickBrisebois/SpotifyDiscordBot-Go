### How to build
Download the repo with go get:

    go get NickBrisebois/SpotifyDiscordBot-Go

cd into the folder:

    cd ~/go/src/NickBrisebois/SpotifyDiscordBot-Go
    
then build it:

    make deps
    make build
    
You will need to have govendor installed which you can install using apt/brew/aur or using go get and having your GOPATH setup correctly. 

### How to run

    ./build/spoticord --config ./config/config.toml
