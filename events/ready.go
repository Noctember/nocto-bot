package events

import (
	"Noctobot/utils"
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	ksoftgo "github.com/KSoft-Si/KSoftgo"
	"github.com/Noctember/hypixel"
	"github.com/jonas747/discordgo"
)

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	logger.Infof("Logged in as %s (%d)", r.User.String(), r.User.ID)
	s.UpdateStatus(0, fmt.Sprintf("=help | %d Servers!", len(s.State.Guilds)))

	utils.Ksoft, _ = ksoftgo.New(utils.GetConfig("ksoft"))
	utils.Hypixel, _ = hypixel.New(utils.GetConfig("hypixel"))
	utils.DBL, _ = dbl.NewClient(utils.GetConfig("dbl"))
	for _, g := range r.Guilds {
		if g.Unavailable {
			guildsCache[g.ID] = true
		}
	}
}
