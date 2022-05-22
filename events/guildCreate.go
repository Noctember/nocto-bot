package events

import (
	"Noctobot/database"
	"Noctobot/utils"
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	"github.com/Noctember/gocto"
	"github.com/dustin/go-humanize"
	"github.com/jonas747/discordgo"
	"github.com/spf13/cast"
	"strconv"
)

func GuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		guildsCache[g.ID] = true
		return
	}

	if _, ok := guildsCache[g.ID]; ok || guildsCache[g.ID] {
		delete(guildsCache, g.ID)

		if len(guildsCache) == 0 {
			logger.Infof("Done loading all guilds (count: %d)", len(s.State.Guilds))
		}

		return
	}

	s.UpdateStatus(0, fmt.Sprintf("=help | %d Servers!", len(s.State.Guilds)))

	owner, err := s.User(g.OwnerID)
	if err != nil {
		return
	}

	created, _ := utils.SnowflakeTimestamp(g.ID)

	_, err = database.GetConfig(g.ID)
	if err == nil {
		database.CreateConfig(g.ID)
	}

	s.ChannelMessageSendEmbed(utils.GuildLogsChannelID, gocto.NewEmbed().
		SetTitle(fmt.Sprintf("%s has joined a new server!", s.State.User.Username)).
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
