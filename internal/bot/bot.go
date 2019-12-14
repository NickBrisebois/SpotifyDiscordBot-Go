package bot

import (
	"errors"
	"fmt"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/config"
	"github.com/NickBrisebois/SpotifyDiscordBot-Go/internal/spotifyauth"
	"github.com/bwmarrin/discordgo"
	"github.com/mvdan/xurls"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	spottyConf *config.Config
	spottyChan chan string
)

// InitBot initializes the discord bot portion of spottybot
func InitBot(config *config.Config) (err error) {
	spottyConf = config

	discord, err := discordgo.New("Bot " + spottyConf.DiscordToken)

	if err != nil {
		fmt.Println("error creating Discord session, ", err)
		err = errors.New("Error initializing bot: " + err.Error())
		return err
	}

	spottyChan = make(chan string)

	// Start spotify API handler
	client := spotifyauth.GetSpotifyClient(spottyConf)
	go InitSpotify(config, spottyChan, client)

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		err = errors.New("Error starting bot: " + err.Error())
		return err
	}

	// Make sure program is killable with signals
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Bot is now running.")

	// If we got here, we received a quit signal so let the spotify thread know that
	spottyChan <- "quit"

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
		urlExtractor := xurls.Relaxed()
		extracted := urlExtractor.FindAllString(m.Content, -1)

		for _, u := range extracted {
			m, err := url.Parse(u)

			if err != nil {
				log.Fatal(err)
			}

			// Make sure the URL we're checking is a spotify URL
			if m.Host == "open.spotify.com" {
				// Grab the ID and ignore the /track/ portion
				var trackPath, ID string
				fmt.Sscanf(m.Path, "%7s%s", &trackPath, &ID)
				spottyChan <- ID
			}
		}
	}

}
