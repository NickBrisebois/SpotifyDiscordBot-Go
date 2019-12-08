package internal

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"syscall"
)

func InitBot(config *Config) (err error) {
	discord, err := discordgo.New("Bot " + config.DiscordToken)

	if err != nil {
		fmt.Println("error creating Discord session, ", err)
		err = errors.New("Error initializing bot: " + err.Error())
		return err
	}

	discord.AddHandler(messageCreate)

	go InitSpotify(config)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		err = errors.New("Error starting bot: " + err.Error())
		return err
	}

	fmt.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == ":)" {
		s.ChannelMessageSend(m.ChannelID, ":(")
	}

}
