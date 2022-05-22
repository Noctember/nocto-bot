package util

import (
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"strings"
)

func Suggest(ctx *gocto.CommandContext) {
	if len(ctx.RawArgs) < 2 {
		ctx.Reply("%s Try to give more details.", utils.CrossMark)
		return
	}

	_, err := ctx.Session.ChannelMessageSendComplex(utils.SuggestionsID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Description: strings.Join(ctx.RawArgs, " "),
			Footer: &discordgo.MessageEmbedFooter{
				Text:    fmt.Sprintf("Guild: %s (%d) | Channel: %s (%d)", ctx.Guild.Name, ctx.Guild.ID, ctx.Channel.Name, ctx.Channel.ID),
				IconURL: discordgo.EndpointGuildIcon(ctx.Guild.ID, ctx.Guild.Icon),
			},
			Author: &discordgo.MessageEmbedAuthor{
				URL:     fmt.Sprintf("https://discordapp.com/channels/%d/%d/%d", ctx.Guild.ID, ctx.Channel.ID, ctx.Message.ID),
				Name:    fmt.Sprintf("%s#%s", ctx.Author.Username, ctx.Author.Discriminator),
				IconURL: ctx.Author.AvatarURL("1024"),
			},
		},
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.Reply("%s Thank you for your suggestion!", utils.CheckMark)
}
