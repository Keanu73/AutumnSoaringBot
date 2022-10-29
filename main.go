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
	bot, err := disgolf.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialise session: %w", err))
	}

	// Adds disgolf command interaction handler & ready event
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(
		func(*discordgo.Session, *discordgo.Ready) {
			log.Println("Autumn Soaring bot is online & ready!")
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
