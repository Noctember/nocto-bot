package info

import (
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/dustin/go-humanize"
	"github.com/jonas747/discordgo"
	"strconv"
)

func Server(ctx *gocto.CommandContext) {
	guild := ctx.Guild
	if len(ctx.RawArgs) > 0 {
		var err error
		guild, err = ctx.Session.Guild(ctx.Arg(0).AsInt64())
		if err != nil {
			ctx.Reply("%s Server not found", utils.CrossMark)
			return
		}
	}
	owner, _ := ctx.FetchUser(guild.OwnerID)
	name, id, region, verification := guild.Name, guild.ID, guild.Region, guild.VerificationLevel
	text, voice, member, human, bot, online, idle, dnd, offline := 0, 0, guild.MemberCount, 0, 0, 0, 0, 0, 0

	for _, cha := range guild.Channels {
		switch cha.Type {
		case discordgo.ChannelTypeGuildNews:
			fallthrough
		case discordgo.ChannelTypeGuildStore:
			fallthrough
		case discordgo.ChannelTypeGuildText:
			text++
		case discordgo.ChannelTypeGuildVoice:
			voice++
		}
	}
	for _, mem := range guild.Members {
		u, _ := ctx.FetchUser(mem.User.ID)
		if u.Bot {
			bot++
		} else {
			human++
		}
	}
	for _, pres := range guild.Presences {
		switch pres.Status {
		case discordgo.StatusOnline:
			online++
		case discordgo.StatusIdle:
			idle++
		case discordgo.StatusDoNotDisturb:
			dnd++
		case discordgo.StatusInvisible:
			offline++
		default:
			fallthrough
		case discordgo.StatusOffline:
			offline++
		}
	}
	created, _ := utils.SnowflakeTimestamp(id)
	if owner == nil {
		owner = &discordgo.User{
			ID:            0,
			Username:      "Unknown",
			Discriminator: "0000",
		}
	}
	embed := gocto.NewEmbed().SetThumbnail(discordgo.EndpointGuildIcon(guild.ID, guild.Icon)).
		SetTitle(name).
		AddInlineField("ID", strconv.FormatInt(id, 10)).
		AddInlineField("Owner", fmt.Sprintf("%s#%s", owner.Username, owner.Discriminator)).
		AddField("Created", fmt.Sprintf("%s (%s)", created.Format("2 January 2006"), humanize.Time(created)))
	if guild.ID == ctx.Guild.ID {
		embed.AddInlineField("Channel", fmt.Sprintf("Text %d | Voice %d", text, voice)).
			AddInlineField("Member", fmt.Sprintf("User %d | Human %d | <:botTag:230105988211015680> %d |"+
				" <:online:313956277808005120> %d | <:away:313956277220802560> %d | <:dnd:313956276893646850> %d | "+
				"<:offline:313956277237710868> %d", member, human, bot, online, idle, dnd, offline))
	}

	embed.AddField("Other", fmt.Sprintf("Region %s | Verification %d", region, verification))

	ctx.ReplyEmbed(embed.Build())
}
