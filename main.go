package main

import (
	"RUSLANGO/events"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {

	token := os.Getenv("RUSLAN_BOT_DISCORD_TOKEN")
	// Create a new Discord session using the provided bot token.
	// log.Fatal exits with an error code, so systemd knows to restart the bot
	// (a plain return exits "successfully", and systemd leaves it stopped)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session, ", err)
	}

	// Register the event handlers. safe() makes sure one broken message can't crash the whole bot.
	dg.AddHandler(safe(events.OnMessage))
	dg.AddHandler(safe(events.OnServerJoin))
	dg.AddHandler(safe(events.OnServerLeave))
	dg.AddHandler(safe(events.OnBotReady))

	// Privileged intents (like server members) need both: the toggle in the developer portal
	// only allows the bot to use them, and this line actually asks Discord for them
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection, ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C like 1000 times in a row to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc //anek is dead...
	// Cleanly close down the Discord session.
	dg.Close()
}

// discordgo runs every handler on its own and doesn't catch panics, so without this
// a single panic (like a bad "roll") takes down the whole bot. This logs it instead.
func safe[T any](handler func(*discordgo.Session, T)) func(*discordgo.Session, T) {
	return func(s *discordgo.Session, event T) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("handler panic:", r)
			}
		}()
		handler(s, event)
	}
}
