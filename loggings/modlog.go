package loggings

import (
	"Noctobot/database"
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"strings"
)

type ModlogAction struct {
	Prefix string
	Emoji  string
	Color  int
	Footer string
}

var (
	MAMute       = ModlogAction{Prefix: "Muted", Emoji: "🔇", Color: 0x57728e}
	MAUnmute     = ModlogAction{Prefix: "Unmuted", Emoji: "🔊", Color: 0x62c65f}
	MAKick       = ModlogAction{Prefix: "Kicked", Emoji: "👢", Color: 0xf2a013}
	MABanned     = ModlogAction{Prefix: "Banned", Emoji: "🔨", Color: 0xd64848}
	MAUnbanned   = ModlogAction{Prefix: "Unbanned", Emoji: "🔓", Color: 0x62c65f}
	MAWarned     = ModlogAction{Prefix: "Warned", Emoji: "⚠", Color: 0xfca253}
	MAGiveRole   = ModlogAction{Prefix: "", Emoji: "➕", Color: 0x53fcf9}
	MARemoveRole = ModlogAction{Prefix: "", Emoji: "➖", Color: 0x53fcf9}
)

func CreateModlogEmbed(ctx *gocto.CommandContext, channelID int64, author *discordgo.User, action ModlogAction, target *discordgo.User, reason string, p *database.Case) error {
	if channelID == 0 {
		return nil
	}

	if author == nil {
		author = &discordgo.User{
			ID:            0,
			Username:      "Unknown",
			Discriminator: "????",
		}
	}

	if reason == "" {
		reason = "No reason provided"
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s#%s (ID %d)", author.Username, author.Discriminator, author.ID),
			IconURL: discordgo.EndpointUserAvatar(author.ID, author.Avatar),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: discordgo.EndpointUserAvatar(target.ID, target.Avatar),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Case #%d", p.ID),
		},
		Color: action.Color,
		Description: fmt.Sprintf("**%s%s %s**#%s *(ID %d)*\n📄**Reason:** %s",
			action.Emoji, action.Prefix, target.Username, target.Discriminator, target.ID, reason),
	}

	_, err := ctx.Session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		ctx.Error(err)
		return err
	}

	return err
}

func SendPunishDM(ctx *gocto.CommandContext, action ModlogAction, target *discordgo.User, reason string) error {
	channel, err := ctx.Session.UserChannelCreate(target.ID)
	if err != nil {
		ctx.Reply("%s %s#%s was %s but did not get a dm due to having DMs off.", utils.CrossMark, target.Username, target.Discriminator, strings.ToLower(action.Prefix))
		return err
	}
	_, err = ctx.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("You have been %s from %s for %s", action.Prefix, ctx.Guild.Name, reason))
	if err != nil {
		ctx.Reply("%s %s#%s was %s but did not get a dm due to having DMs off.", utils.CrossMark, target.Username, target.Discriminator, strings.ToLower(action.Prefix))
		return err
	}
	return nil
}
