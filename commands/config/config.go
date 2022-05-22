package config

import (
	"Noctobot/database"
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"strings"
)

func Config(ctx *gocto.CommandContext) {
	key := ctx.Arg(0).AsString()
	val := strings.Join(ctx.RawArgs[1:], " ")
	database.GetGConfig(ctx.Guild.ID)
	err := database.SetValue(ctx.Guild.ID, key, val)
	if err != nil {
		ctx.Error(err)
		ctx.Reply("%s Config key `%s` does not exist. Available keys: `qrcode`, `case`", utils.CrossMark, key)
		return
	} else {
		ctx.Reply("%s Changed `%s` to `%s`", utils.CheckMark, key, val)
	}
}
