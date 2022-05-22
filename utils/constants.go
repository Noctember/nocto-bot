package utils

import (
	"github.com/DiscordBotList/go-dbl"
	ksoftgo "github.com/KSoft-Si/KSoftgo"
	"github.com/Noctember/hypixel"
)

const (
	OwnerID            = 118499178790780929
	LogsChannelID      = 346897238066331649
	GuildLogsChannelID = 346897238066331649
	SuggestionsID      = 673042654493278229
	BotID              = 328239915965874177
	CrossMark          = "❌"
	CheckMark          = "✅"
	Version            = "1.0.0"
)

var (
	Ksoft   *ksoftgo.KSession
	Hypixel *hypixel.Hypixel
	DBL     *dbl.DBLClient
)
