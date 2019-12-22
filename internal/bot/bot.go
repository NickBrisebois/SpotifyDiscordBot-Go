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
	"strings"
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
	spottyHandler := NewSpotifyHandler(config, spottyChan, client)
	go spottyHandler.SpottyLoop()

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
		// TODO: Redo these commands into a better scaleable way
		if m.Content == "!playlist" {
			link := "https://open.spotify.com/playlist/" + spottyConf.SpotifyPlaylist
			s.ChannelMessageSend(m.ChannelID, link)
		} else if strings.HasPrefix(m.Content, "!whoadded") {
			spotifyIds := getSpotifyID(m.Content)
			for _, id := range spotifyIds {
				addedByRequest := types.SpottyMessage{
					Kind: types.AddedBy,
					ID:   id,
				}
				spottyChan <- addedByRequest
				// Wait for reply
				spottyResp := <-spottyChan

				if spottyResp.Err != nil {
					log.Fatal(spottyResp.Err)
				} else {
					s.ChannelMessageSend(m.ChannelID, spottyResp.User)
				}
			}
		} else {
			spotifyIds := getSpotifyID(m.Content)
			for _, id := range spotifyIds {
				response := handleTrack(id, m.Author.Username)
				s.ChannelMessageSend(m.ChannelID, response)
			}
		}
	}
}

func getSpotifyID(msg string) []string {
	urlExtractor := xurls.Relaxed()
	extracted := urlExtractor.FindAllString(msg, -1)

	var ids []string
	for _, u := range extracted {
		u, err := url.Parse(u)

		if err != nil {
			log.Fatal(err)
		}

		if u.Host == "open.spotify.com" {
			var trackPath, ID string
			fmt.Sscanf(u.Path, "%7s%s", &trackPath, &ID)

			if trackPath == "/track/" {
				ids = append(ids, ID)
			}
		}
	}

	return ids
}

func handleTrack(ID string, user string) (response string) {
	// Send the spotify ID to the spotify API handling thread
	spottyMsg := types.SpottyMessage{
		Kind: types.Track,
		User: user,
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
