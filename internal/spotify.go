package internal

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	auth  spotify.Authenticator
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

// InitSpotify starts spotify API handler
func InitSpotify(config *Config) (err error) {
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

	const PlaylistID spotify.ID = "37i9dQZF1EthtctLd3ak1i"
	results, err := client.GetPlaylist(PlaylistID)
	if err != nil {
		log.Fatalf("couldn't get features playlists: %v", err)
	}

	fmt.Println(results.Name)

	// Listen for signals so that the program is capable of quitting
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}
