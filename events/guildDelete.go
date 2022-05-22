package events

import (
	"Noctobot/utils"
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	"github.com/Noctember/gocto"
	"github.com/dustin/go-humanize"
	"github.com/jonas747/discordgo"
	"github.com/spf13/cast"
	"strconv"
)

func GuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	if g.Unavailable {
		guildsCache[g.ID] = true
		return
	}
	s.UpdateStatus(0, fmt.Sprintf("=help | %d Servers!", len(s.State.Guilds)))

	created, _ := utils.SnowflakeTimestamp(g.ID)
	owner, _ := s.User(g.OwnerID)

	s.ChannelMessageSendEmbed(utils.GuildLogsChannelID, gocto.NewEmbed().
		SetTitle(fmt.Sprintf("%s has left a new server!", s.State.User.Username)).
		SetDescription(g.Name).
		SetThumbnail(discordgo.EndpointGuildIcon(g.ID, g.Icon)).
		SetColor(0xDFAC7C).
		AddField("Owner", owner.String()).
		AddField("Member Count", strconv.Itoa(g.MemberCount)).
		AddField("Created At", fmt.Sprintf("%s (%s)", created.Format("2 January 2006"), humanize.Time(created))).
		SetFooter(strconv.FormatInt(g.ID, 10)).
		Build())

	if !cast.ToBool(utils.GetConfig("dev")) && utils.DBL != nil {
		utils.DBL.PostBotStats(strconv.FormatInt(s.State.User.ID, 10), dbl.BotStatsPayload{
			Shards: []int{len(s.State.Guilds)},
		})
	}
}
