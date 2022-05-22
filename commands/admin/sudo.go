package admin

import "github.com/Noctember/gocto"

func OwnerSu(ctx *gocto.CommandContext) {
	cmd := ctx.Bot.GetCommand(ctx.Arg(1).AsString())

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
		Author:      ctx.Arg(0).AsUser(),
		RawArgs:     ctx.RawArgs[2:],
		Prefix:      ctx.Prefix,
		Guild:       ctx.Guild,
		Flags:       ctx.Flags,
		Locale:      ctx.Locale,
		InvokedName: ctx.Arg(1).AsString(),
	}

	if !cctx.ParseArgs() {
		return
	}

	cmd.Run(cctx)
}

func OwnerSudo(ctx *gocto.CommandContext) {
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
		Prefix:      ctx.Prefix,
		Guild:       ctx.Guild,
		Flags:       ctx.Flags,
		Locale:      ctx.Locale,
		InvokedName: ctx.Arg(0).AsString(),
	}

	if !cctx.ParseArgs() {
		return
	}

	cmd.Run(cctx)
}
