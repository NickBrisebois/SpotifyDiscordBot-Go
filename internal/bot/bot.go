package bot

import (
	"errors"
	"fmt"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/spotifyauth"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/types"
	"github.com/bwmarrin/discordgo"
	"github.com/mvdan/xurls"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	spottyConf *types.Config
	spottyChan chan types.SpottyMessage
	spottyResp string
)

// InitBot initializes the discord bot portion of spottybot
func InitBot(config *types.Config) (err error) {
	spottyConf = config

	discord, err := discordgo.New("Bot " + spottyConf.DiscordToken)

	if err != nil {
		log.Println("error creating Discord session, ", err)
		err = errors.New("Error initializing bot: " + err.Error())
		return err
	}

	spottyChan = make(chan types.SpottyMessage)

	// Start spotify API handler
	client := spotifyauth.GetSpotifyClient(spottyConf)
	go InitSpotify(config, spottyChan, client)

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		err = errors.New("Error starting bot: " + err.Error())
		return err
	}

	log.Println("Bot is now running.")

	// Make sure program is killable with signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// If we got here, we received a quit signal so let the spotify thread know that
	spottyMsg := types.SpottyMessage{
		Msg: "quit",
	}
	spottyChan <- spottyMsg

	discord.Close()
	return nil
}

// messageCreate handles when a message has been sent and should be responded to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Skip messages by the bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	if spottyConf.LimitToOneChannel == false || m.ChannelID == spottyConf.ChannelToUse {
		//		s.ChannelMessageSend(m.ChannelID, m.Content)
		if m.Content == "!playlist" {
			link := "https://open.spotify.com/playlist/" + spottyConf.SpotifyPlaylist
			s.ChannelMessageSend(m.ChannelID, link)
		} else {
			urlExtractor := xurls.Relaxed()
			extracted := urlExtractor.FindAllString(m.Content, -1)

			for _, u := range extracted {
				u, err := url.Parse(u)

				if err != nil {
					log.Fatal(err)
				}

				// Make sure the URL we're checking is a spotify URL
				if u.Host == "open.spotify.com" {
					// Grab the ID and ignore the /track/ portion
					var trackPath, ID string
					fmt.Sscanf(u.Path, "%7s%s", &trackPath, &ID)

					if trackPath == "/track/" {
						response := handleTrack(ID)
						s.ChannelMessageSend(m.ChannelID, response)
					} else if trackPath == "/album/" {
						s.ChannelMessageSend(m.ChannelID, "Should adding albums even be a feature?")
					} else if trackPath == "/playlist" {
						s.ChannelMessageSend(m.ChannelID, "Should adding playlists even be a feature?")
					}
				}
			}
		}
	}
}

func handleTrack(ID string) (response string) {
	// Send the spotify ID to the spotify API handling thread
	spottyMsg := types.SpottyMessage{
		Kind: types.Track,
		ID:   ID,
	}
	spottyChan <- spottyMsg
	// Wait for reply
	spottyResp := <-spottyChan

	if spottyResp.Err != nil {
		var errorMessage string
		errorMessage = "Look what you made me do: " + spottyResp.Err.Error()
		return errorMessage
	}

	return spottyResp.Msg
}
