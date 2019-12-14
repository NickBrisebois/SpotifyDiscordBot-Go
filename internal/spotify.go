package internal

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
)

var (
	auth  spotify.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// InitSpotify starts spotify API handler
func InitSpotify(config *Config, spottyChan chan string) (err error) {
	spotifyConfig := &clientcredentials.Config{
		ClientID:     config.SpotifyClientID,
		ClientSecret: config.SpotifyClientSecret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := spotifyConfig.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	client := spotify.Authenticator{}.NewClient(token)

	var playlistID spotify.ID
	playlistID = spotify.ID(spottyConf.SpotifyPlaylist)
	results, err := client.GetPlaylist(playlistID)
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(results.Name)

	for {
		data := <-spottyChan
		fmt.Println(data)

		if data == "exit" {
			break
		}
	}

	// Listen for signals so that the program is capable of quitting
	//sc := make(chan os.Signal, 1)
	//signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	//<-sc

	return nil
}
