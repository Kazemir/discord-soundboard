package main

import (
	"fmt"
	"os"

	"github.com/Kazemir/discord-soundboard/feature/window"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	discordBotToken  = os.Getenv("DISCORD_BOT_TOKEN")
	discordGuildId   = os.Getenv("DISCORD_GUILD_ID")
	discordChannelId = os.Getenv("DISCORD_CHANNEL_ID")
	isDebug          = os.Getenv("IS_DEBUG") == "true"
)

func main() {
	if discordBotToken == "" {
		fmt.Println("DISCORD_BOT_TOKEN is not set")
		return
	}
	if discordGuildId == "" {
		fmt.Println("DISCORD_GUILD_ID is not set")
		return
	}
	if discordChannelId == "" {
		fmt.Println("DISCORD_CHANNEL_ID is not set")
		return
	}

	if isDebug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	discord, err := discordgo.New("Bot " + discordBotToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Error creating Discord session!")
		return
	}

	discord.Identify.Intents |= discordgo.IntentsGuilds
	discord.Identify.Intents |= discordgo.IntentsGuildVoiceStates
	discord.Identify.Intents |= discordgo.IntentsGuildPresences
	discord.Identify.Intents |= discordgo.IntentsGuildMessages

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info("Discord bot is ready!")
	})

	discord.Open()

	w := window.CreateWindow(discord, discordGuildId, discordChannelId)
	w.SetOnClosed(func() {
		discord.Close()
	})
	w.ShowAndRun()
}
