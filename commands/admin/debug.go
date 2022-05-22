package admin

import (
	"fmt"
	"github.com/Noctember/gocto"
	"time"
)

func Debug(ctx *gocto.CommandContext) {
	before := time.Now()
	cmd := ctx.Bot.GetCommand(ctx.Arg(0).AsString())

	if cmd == nil {
		ctx.Reply("Invalid command.")
		return
	}

	cctx := &gocto.CommandContext{
		Bot:         ctx.Bot,
		Command:     cmd,
		Message:     ctx.Message,
		Channel:     ctx.Channel,
		Session:     ctx.Session,
		Author:      ctx.Author,
		RawArgs:     ctx.RawArgs[1:],
		Args:        ctx.Args[1:],
		Prefix:      ctx.Prefix,
		Guild:       ctx.Guild,
		Flags:       ctx.Flags,
		Locale:      ctx.Locale,
		InvokedName: ctx.InvokedName,
	}

	if !cctx.ParseArgs() {
		return
	}

	cmd.Run(cctx)

	after := time.Now()
	ti := after.Sub(before).Milliseconds()
	ctx.Reply(fmt.Sprintf("`%s` took %d ms to complete", cmd.Name, ti))
}
