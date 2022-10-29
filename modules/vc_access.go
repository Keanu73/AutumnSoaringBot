package modules

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Keanu73/AutumnSoaringBot/cron"
	"github.com/bwmarrin/discordgo"
)

var accessTimes = []struct {
	startTime string
	lock      bool
}{
	{"2022-10-30 06:25:00", false},
	{"2022-10-30 08:00:00", true},
}

func (*vcaccess) ScheduleCrons(session *discordgo.Session) {
	guildID := os.Getenv("BOT_GUILD_ID")
	vcID := os.Getenv("BOT_VC_ID")

	// Loops through access times & schedules cron for them
	for _, job := range accessTimes {
		lockStr := ""
		if job.lock {
			lockStr = "lock"
		} else {
			lockStr = "unlock"
		}

		log.Printf("[VC ACCESS] Scheduling cron to %s the VC at: %s", lockStr, job.startTime)
		job := job
		go func() {
			scheduleAccessModification(job.startTime, session, guildID, vcID, job.lock)
		}()
	}
}

func scheduleAccessModification(
	startTime string, session *discordgo.Session, guildID string, channelID string,
	lock bool,
) {
	ctx := context.Background()

	start, err := time.ParseInLocation(
		"2006-01-02 15:04:05",
		startTime,
		time.Local,
	) // any time in the past works, but it should be on the hour
	// uses local TZ - in this case, GMT (set in Dockerfile)
	if err != nil {
		log.Println(fmt.Errorf("[VC ACCESS] unable to parse start time: %w", err))
		return
	}

	interval := time.Hour * 24 // 24 hours

	for range cron.Every(ctx, start, interval) {
		lockStr := ""
		if lock {
			lockStr = "Locking"
		} else {
			lockStr = "Unlocking"
		}

		log.Printf("[VC ACCESS]: %s at %s", lockStr, startTime)

		err = modifyAccess(session, guildID, channelID, lock)
		if err != nil {
			log.Println(err)
		}
	}
}

func modifyAccess(session *discordgo.Session, guildID string, channelID string, lock bool) error {
	// basically, deny/allow @everyone the permission to connect to the VC
	// requires that the bot has the "Connect" permission allowed for itself, otherwise it cannot modify permission

	var allowPermissionInt int64
	var denyPermissionInt int64

	if lock {
		// permission int for "Connect" to VC
		denyPermissionInt = 1048576
	} else {
		allowPermissionInt = 1048576
	}

	// we're using guild ID as target ID to specify @everyone
	err := session.ChannelPermissionSet(
		channelID, guildID, discordgo.PermissionOverwriteTypeRole, allowPermissionInt,
		denyPermissionInt,
	)

	if err != nil {
		return fmt.Errorf("[VC ACCESS] unable to set permissions on channel: %w", err)
	}

	return nil
}
