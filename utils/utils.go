package utils

import (
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"github.com/pollen5/minori"
	"time"
)

var Logger = minori.GetLogger("Noctobot")

func ErrorHandler(bot *gocto.Bot, err interface{}) {
	if cmd, ok := err.(*gocto.CommandError); ok {
		ctx := cmd.Context

		g := "DM"
		gid := ctx.Channel.ID

		if ctx.Guild != nil {
			g = ctx.Guild.Name
			gid = ctx.Guild.ID
		}
		bot.Session.ChannelMessageSendEmbed(LogsChannelID, &discordgo.MessageEmbed{
			Title:       "Command Error - " + ctx.Command.Name,
			Description: fmt.Sprintf("Content: %s\n```\n%s```", ctx.Message.Content, cmd.Error()),
			Color:       0xDFAC7C,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Author: %s (%d), Guild: %s (%d)", ctx.Author.Username, ctx.Author.ID, g, gid),
			},
		})

		Logger.Errorf("Command Error (%s, %s:%d): %s (Author: %s (%d), Guild: %s (%d))", ctx.Command.Name, cmd.File, cmd.Line, cmd.Error(), ctx.Author.Username, ctx.Author.ID, g, gid)
		ctx.ReplyLocale("COMMAND_ERROR")
		return
	}

	bot.Session.ChannelMessageSendEmbed(LogsChannelID, &discordgo.MessageEmbed{
		Title:       "Panic Recovered",
		Description: fmt.Sprint(err),
		Color:       0xDFAC7C,
	})

	Logger.Error(err)
}

func SnowflakeTimestamp(ID int64) (t time.Time, err error) {
	timestamp := (ID >> 22) + 1420070400000
	t = time.Unix(timestamp/1000, 0)
	return
}
