package util

import (
	"fmt"
	"github.com/Noctember/gocto"
)

func Avatar(ctx *gocto.CommandContext) {
	user := ctx.Author

	if ctx.HasArgs() {
		user = ctx.Arg(0).AsUser()
	}

	ctx.BuildEmbed(gocto.NewEmbed().
		SetTitle(fmt.Sprintf("%s's Avatar", user.Username)).
		SetDescription(fmt.Sprintf("Sizes: [128](%s) | [512](%s) | [1024](%s)", user.AvatarURL("128"), user.AvatarURL("512"), user.AvatarURL("1024"))).
		SetColor(0xDFAC7C).
		SetImage(user.AvatarURL("2048")).
		SetAuthor(ctx.Author.Username, ctx.Author.AvatarURL("256")))
}
