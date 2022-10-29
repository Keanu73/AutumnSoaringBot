package modules

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Keanu73/AutumnSoaringBot/cron"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/pkg/errors"
)

/*
we want to:
 * schedule modules crons.
 * automatically unlock VC between 6:25-8am.
 how?

 * 1) use mongo collection - "modules" to store time intervals & their associated mp3 files.
 * 2) query mongo collection
 * 3) schedule cron for each time, & play mp3 in VC at specified time.
*/

// We're going to use British times here
// The Europe/London TZ WILL be enforced in the Dockerfile.
var filesToPlay = []struct {
	startTime string
	fileName  string
}{
	{"2022-10-28 06:30:00", "workout.dca"},
	{"2022-10-28 06:40:00", "meditate.dca"},
	{"2022-10-28 06:45:00", "journal.dca"},
	{"2022-10-28 06:50:00", "discussion.dca"},
}

func (*timekeeping) ScheduleCrons(session *discordgo.Session) {
	guildID := os.Getenv("BOT_GUILD_ID")
	vcID := os.Getenv("BOT_VC_ID")

	for _, file := range filesToPlay {
		log.Printf("[TIMEKEEPING] Scheduling file: %s %s %s %s", file.startTime, file.fileName, guildID, vcID)
		file := file
		// Uses goroutines to schedule the crons in parallel
		go func() {
			scheduleFile(file.startTime, session, guildID, vcID, file.fileName)
		}()
	}
}

func scheduleFile(
	startTime string, session *discordgo.Session, guildID string, channelID string,
	fileName string,
) {
	ctx := context.Background()

	start, err := time.ParseInLocation(
		"2006-01-02 15:04:05",
		startTime,
		time.Local,
	) // any time in the past works, but it should be on the hour
	// Uses time.Local for TZ of environment (Docker in this case)
	if err != nil {
		log.Println(fmt.Errorf("[TIMEKEEPING] unable to parse start time: %w", err))
		return
	}

	interval := time.Hour * 24 // 24 hours

	for range cron.Every(ctx, start, interval) {
		// log.Printf("Running file: %s %s %s %s", startTime, fileName, guildID, channelID)
		log.Printf("[TIMEKEEPING]: Playing %s", fileName)

		err = playFile(session, guildID, channelID, fileName)
		if err != nil {
			log.Println(err)
		}
	}
}

func playFile(session *discordgo.Session, guildID string, channelID string, fileName string) error {
	// Connect to VC
	voiceConnection, err := session.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return fmt.Errorf("[TIMEKEEPING] couldn't join voice channel: %w", err)
	}

	file, err := os.ReadFile(fmt.Sprintf("./audio/dca/%s", fileName))
	if err != nil {
		return fmt.Errorf("[TIMEKEEPING] couldn't read DCA file: %w", err)
	}

	inputReader := bytes.NewReader(file)

	// inputReader is an io.Reader, like a file for example
	decoder := dca.NewDecoder(inputReader)

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return fmt.Errorf("[TIMEKEEPING] couldn't read file: %w", err)
			}

			break
		}

		// Do something with the frame, in this example were sending it to discord
		select {
		case voiceConnection.OpusSend <- frame:
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			return errors.New("[TIMEKEEPING] connection lost")
		}
	}
	voiceConnection.Close()

	// Leave VC
	_, _ = session.ChannelVoiceJoin(guildID, "", false, false)

	return nil
}
