package mod

import (
	"Noctobot/database"
	"Noctobot/loggings"
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
	"strings"
)

func Warns(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	cases := database.GetCases(target.ID, ctx.Guild.ID)
	if len(cases) < 1 {
		ctx.Reply("%s This user hasn't been warned yet.", utils.CrossMark)
		return
	}
	p := gocto.NewPaginatorForContext(ctx)
	p.Delete()
	p.SetExtra(fmt.Sprintf("%d Warns", len(cases)))
	tmp := ""
	for i, c := range cases {
		if c.Type == "warn" {
			tmp += fmt.Sprintf("ID: %d\nMod: %s\nReason: %s\n-----\n", c.ID, c.AuthorUsernameDiscrim, c.Message)
		}
		if i != 0 && i%7 == 0 {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.AddField("Warnings", tmp)
			})
			tmp = ""
		}
	}
	if tmp != "" {
		p.AddPage(func(em *gocto.Embed) *gocto.Embed {
			return em.AddField("Warnings", tmp)
		})
	}

	go p.Run()
}

func DeleteCase(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsInt64()

	_, err := database.Exec("DELETE FROM modlogs WHERE id = $1 AND gid = $2", target, ctx.Guild.ID)
	if err != nil {
		ctx.Reply("%s No case with id `%d` found for `%s`", utils.CrossMark, target)
		return
	}
	ctx.Reply("%s Case deleted.", utils.CheckMark)
	return

}

func Clearwarns(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	c := database.GetCases(target.ID, ctx.Guild.ID)
	if len(c) < 1 {
		ctx.Reply("%s This user has not been warned yet.", utils.CrossMark)
		return
	}

	_, err := database.Exec("DELETE FROM modlogs WHERE uid = $1 AND gid = $2", target.ID, ctx.Guild.ID)
	if err != nil {
		ctx.Reply("%s Could not clear warns.", utils.CrossMark)
	} else {
		ctx.Reply("%s `%s`'s warns have been cleared.", utils.CheckMark, target.Username)
	}
}

func Warn(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsUser()
	reason := strings.Join(ctx.RawArgs[1:], " ")
	warning := &database.Case{
		GuildID:               ctx.Guild.ID,
		UserID:                cast.ToString(target.ID),
		ModID:                 ctx.Author.ID,
		AuthorUsernameDiscrim: ctx.Author.Username + "#" + ctx.Author.Discriminator,
		Type:                  "warn",
		Message:               reason,
	}

	gconfig, _ := database.GetGConfig(ctx.Guild.ID)

	err := database.CreateCase(warning)
	if err == nil {
		loggings.SendPunishDM(ctx, loggings.MAWarned, target, reason)
		loggings.CreateModlogEmbed(ctx, gconfig.CaseChannel, ctx.Author, loggings.MAWarned, target, reason, warning)
	} else {
		ctx.Error(err)
	}
}
