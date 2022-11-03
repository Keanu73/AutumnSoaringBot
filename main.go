package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FedorLap2006/disgolf"
	"github.com/Keanu73/AutumnSoaringBot/modules"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Checks if env variables were set
	botToken := os.Getenv("BOT_TOKEN")
	guildID := os.Getenv("BOT_GUILD_ID")
	vcID := os.Getenv("BOT_VC_ID")
	mongoURI := os.Getenv("MONGODB_URI")

	if botToken == "" || guildID == "" || vcID == "" || mongoURI == "" {
		log.Fatal(
			"[ERROR] One or more environment variables were not supplied.\n" +
				"Please ensure the following environment variables are supplied:\n" +
				"* BOT_TOKEN\n" +
				"* BOT_GUILD_ID\n" +
				"* BOT_VC_ID\n" +
				"* MONGODB_URI",
		)
	}

	// If all good... let's continue

	bot, err := disgolf.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialise session: %w", err))
	}

	// Adds disgolf command interaction handler & ready event
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(
		func(session *discordgo.Session, ready *discordgo.Ready) {
			log.Printf("[READY] %s#%s is online!", ready.User.Username, ready.User.Discriminator)
		},
	)

	// Opens bot's session using token
	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open session: %w", err))
	}

	// Syncs disgolf router with discordgo bot
	err = bot.Router.Sync(bot.Session, "", os.Getenv("BOT_GUILD_ID"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to sync commands: %w", err))
	}

	// Updates bot status
	_ = bot.Session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "people better themselves",
					Type: discordgo.ActivityTypeWatching,
				},
			},
		},
	)

	// Schedules timekeeping & VC access crons
	modules.Timekeeping.ScheduleCrons(bot.Session)
	modules.VCAccess.ScheduleCrons(bot.Session)

	// Allows for graceful Ctrl + C
	ech := make(chan os.Signal)
	signal.Notify(ech, os.Kill, syscall.SIGTERM) //nolint:govet,staticcheck
	<-ech
	_ = bot.Close()
}
