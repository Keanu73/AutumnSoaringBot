package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/FedorLap2006/disgolf"
	"github.com/Keanu73/AutumnSoaringBot/timekeeping"
	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

func init() {
	err := dotenv.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load .env: %w", err))
	}
}

func main() {
	bot, err := disgolf.New(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialise session: %w", err))
	}
	bot.AddHandler(bot.Router.HandleInteraction)
	bot.AddHandler(func(*discordgo.Session, *discordgo.Ready) { log.Println("Ready!") })
	err = bot.Open()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open session: %w", err))
	}
	err = bot.Router.Sync(bot.Session, "", os.Getenv("BOT_GUILD_ID"))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to sync commands: %w", err))
	}

	// Updates bot status
	_ = bot.Session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: "people become better by the day",
					Type: discordgo.ActivityTypeWatching,
				},
			},
		},
	)

	timekeeping.ScheduleCrons(bot.Session)

	ech := make(chan os.Signal)
	signal.Notify(ech, os.Kill, syscall.SIGTERM) //nolint:govet,staticcheck
	<-ech
	_ = bot.Close()
}
