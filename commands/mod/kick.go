package mod

import (
	"Noctobot/database"
	"Noctobot/loggings"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
	"strings"
)

func Kick(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	reason := strings.Join(ctx.RawArgs[1:], " ")
	kick := &database.Case{
		GuildID:               ctx.Guild.ID,
		UserID:                cast.ToString(target.ID),
		ModID:                 ctx.Author.ID,
		AuthorUsernameDiscrim: ctx.Author.Username + "#" + ctx.Author.Discriminator,
		Type:                  "kick",
		Message:               reason,
	}

	gconfig, _ := database.GetConfig(ctx.Guild.ID)

	err := database.CreateCase(kick)
	if err == nil {
		loggings.SendPunishDM(ctx, loggings.MAKick, target, reason)
		loggings.CreateModlogEmbed(ctx, gconfig.CaseChannel, ctx.Author, loggings.MAKick, target, reason, kick)
		ctx.Session.GuildMemberDeleteWithReason(ctx.Guild.ID, target.ID, "Kicked by "+kick.AuthorUsernameDiscrim+" ("+reason+")")
	} else {
		ctx.Error(err)
	}
}
