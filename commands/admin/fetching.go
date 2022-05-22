package admin

import (
	"encoding/json"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
)

func OwnerFetch(ctx *gocto.CommandContext) {
	switch ctx.Arg(0).AsString() {
	case "g":
		fallthrough
	case "guild":
		OwnerFetchGuild(ctx)
	case "m":
		fallthrough
	case "mem":
		fallthrough
	case "member":
		OwnerFetchMember(ctx)
	case "u":
		fallthrough
	case "user":
		OwnerFetchUser(ctx)
	case "c":
		fallthrough
	case "chan":
		fallthrough
	case "channel":
		OwnerFetchChannel(ctx)
	case "p":
		fallthrough
	case "pres":
		fallthrough
	case "presence":
		OwnerFetchPresence(ctx)
	}
}

func OwnerFetchGuild(ctx *gocto.CommandContext) {
	id := ctx.Arg(1).AsInt64()
	guild, _ := ctx.Session.Guild(cast.ToInt64(id))
	out, _ := json.MarshalIndent(guild, "", "\t")
	ctx.CodeBlock("json", string(out))
}
func OwnerFetchMember(ctx *gocto.CommandContext) {
	id := ctx.Arg(1).AsInt64()
	mem := ctx.Member(cast.ToInt64(id))
	out, _ := json.MarshalIndent(mem, "", "\t")
	ctx.CodeBlock("json", string(out))
}
func OwnerFetchUser(ctx *gocto.CommandContext) {
	id := ctx.Arg(1).AsInt64()
	usr := ctx.User(cast.ToInt64(id))
	out, _ := json.MarshalIndent(usr, "", "\t")
	ctx.CodeBlock("json", string(out))
}
func OwnerFetchChannel(ctx *gocto.CommandContext) {
	id := ctx.Arg(1).AsInt64()
	chann, _ := ctx.Session.Channel(cast.ToInt64(id))
	out, _ := json.MarshalIndent(chann, "", "\t")
	ctx.CodeBlock("json", string(out))
}
func OwnerFetchPresence(ctx *gocto.CommandContext) {
	id := ctx.Arg(1).AsInt64()
	ctx.Member(id)
	chann, _ := ctx.Session.State.Presence(ctx.Guild.ID, cast.ToInt64(id))
	out, _ := json.MarshalIndent(chann, "", "\t")
	ctx.CodeBlock("json", string(out))
}
