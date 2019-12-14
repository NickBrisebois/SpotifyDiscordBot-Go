package bot

import (
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/config"
	"github.com/zmb3/spotify"
	"log"
)

// InitSpotify starts spotify API handler
func InitSpotify(config *config.Config, spottyChan chan string, client *spotify.Client) (err error) {
	var playlistID spotify.ID
	playlistID = spotify.ID(config.SpotifyPlaylist)

	if err != nil {
		log.Fatalf("couldn't get features playlists:%v", err)
	}

	// We will loop every time we get some IDs or the exit command from the spotty channel
	for {
		data := <-spottyChan

		if data == "quit" {
			break
		} else {
			log.Println("Adding track to playlist:", data)
			snapID, err := client.AddTracksToPlaylist(playlistID, spotify.ID(data))
			if err != nil {
				log.Println(err)
			} else {
				log.Println("snap id:", snapID)
			}
		}
	}

	return nil
}
