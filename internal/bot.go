package internal

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func InitBot(config *Config) (err error) {
	fmt.Printf("Discord Token: %s", config.DiscordToken)
	discord, err := discordgo.New("Bot " + config.DiscordToken)

	if err != nil {
		fmt.Println("error creating Discord session, ", err)
		err = errors.New("Error initializing bot: " + err.Error())
		return err
	}

	discord.AddHandler(messageCreate)

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

}
