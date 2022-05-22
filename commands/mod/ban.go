package mod

import (
	"Noctobot/database"
	"Noctobot/loggings"
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
	"strings"
)

func Ban(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	reason := strings.Join(ctx.RawArgs[1:], " ")
	ban := &database.Case{
		GuildID:               ctx.Guild.ID,
		UserID:                cast.ToString(target.ID),
		ModID:                 ctx.Author.ID,
		AuthorUsernameDiscrim: ctx.Author.Username + "#" + ctx.Author.Discriminator,
		Type:                  "ban",
		Message:               reason,
	}

	gconfig, _ := database.GetConfig(ctx.Guild.ID)

	days := 0
	if ctx.HasFlag("prune") {
		days = cast.ToInt(ctx.Flag("prune"))
	}

	err := database.CreateCase(ban)
	if err == nil {
		ctx.Session.GuildBanCreateWithReason(ctx.Guild.ID, target.ID, "Banned by "+ban.AuthorUsernameDiscrim+"(Reason: "+reason+")", days)
		loggings.SendPunishDM(ctx, loggings.MABanned, target, reason)
		loggings.CreateModlogEmbed(ctx, gconfig.CaseChannel, ctx.Author, loggings.MABanned, target, reason, ban)
	} else {
		ctx.Error(err)
	}
}

func Unban(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	unban := &database.Case{
		GuildID:               ctx.Guild.ID,
		UserID:                cast.ToString(target.ID),
		ModID:                 ctx.Author.ID,
		AuthorUsernameDiscrim: ctx.Author.Username + "#" + ctx.Author.Discriminator,
		Message:               "Manual unban",
	}

	gconfig, _ := database.GetConfig(ctx.Guild.ID)
	ctx.Session.GuildBans(ctx.Guild.ID)
	err := ctx.Session.GuildBanDelete(ctx.Guild.ID, target.ID)
	if err != nil {
		loggings.CreateModlogEmbed(ctx, gconfig.CaseChannel, ctx.Author, loggings.MAUnbanned, target, unban.Message, unban)
	} else {
		ctx.Reply("%s This user in not banned.", utils.CrossMark)
	}
}
