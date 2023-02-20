package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(onReady)

	dg.AddHandler(simpleHandler)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	// Wait for a SIGINT or SIGTERM signal to gracefully exit the bot
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	println("Bot is now running. Press CTRL-C to exit.")
}

func simpleHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!yeah" {
		if _, err := s.ChannelMessageSend(m.ChannelID, "BAD TO THE BONE 💀🔥"); err != nil {
			fmt.Println("ERROR: ", err)
		}

		voiceChannel, err := findUserVoiceChannel(s, m.Author.ID)
		if err != nil {
			return
		}

		voiceConnection, err := s.ChannelVoiceJoin(voiceChannel.GuildID, voiceChannel.ID, false, true)
		if err != nil {
			panic(err)
		}

		voiceConnection.Speaking(true)

		dgvoice.PlayAudioFile(voiceConnection, "badtothebone.mp3", make(chan bool))
		voiceConnection.Disconnect()
	}
}

func findUserVoiceChannel(session *discordgo.Session, userID string) (*discordgo.Channel, error) {
	user, err := session.User(userID)
	if err != nil {
		return nil, err
	}

	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		return nil, err
	}

	for _, guild := range guilds {
		state, err := session.State.Guild(guild.ID)
		if err != nil {
			continue
		}

		for _, vs := range state.VoiceStates {
			if vs.UserID == user.ID {
				channel, err := session.Channel(vs.ChannelID)
				if err != nil {
					return nil, err
				}
				return channel, nil
			}
		}
	}

	return nil, errors.New("User is not in a voice channel")
}
