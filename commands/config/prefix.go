package config

import (
	"Noctobot/database"
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"strings"
)

func SetPrefix(ctx *gocto.CommandContext) {
	if !ctx.HasArgs() {
		ctx.Reply("%s The current prefix is `%s`", utils.CheckMark, database.GetPrefix(ctx.Guild.ID))
		return
	}

	prefix := strings.ReplaceAll(ctx.JoinedArgs(), `"`, ``)

	if len(prefix) > 10 {
		ctx.Reply("%s Prefix must be shorter than 10 characters.", utils.CrossMark)
		return
	}

	_, err := database.Exec(`INSERT INTO guilds (id, prefix, "case") VALUES ($1, $2, '') ON CONFLICT (id) DO UPDATE SET prefix=$2`, ctx.Guild.ID, prefix)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Reply("%s Successfully updated the prefix for this server to `%s`", utils.CheckMark, prefix)
}

func Prefix(ctx *gocto.CommandContext) {
	ctx.Reply("%s The current prefix is `%s`", utils.CheckMark, database.GetPrefix(ctx.Guild.ID))
}
