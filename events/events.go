package events

import (
	"github.com/jonas747/discordgo"
	"github.com/pollen5/minori"
)

var logger = minori.GetLogger("Events")

var guildsCache = make(map[int64]bool)

func Init(s *discordgo.Session) {
	s.AddHandler(Ready)
	s.AddHandler(GuildCreate)
	s.AddHandler(GuildDelete)
}
