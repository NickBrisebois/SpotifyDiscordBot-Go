package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/bot"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/types"
	"log"
)

func main() {
	configPath := flag.String("config", "./config.toml", "Path to config.toml")
	flag.Parse()

	var config types.Config
	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		log.Fatal(err)
	}

	bot.InitBot(&config)
}
