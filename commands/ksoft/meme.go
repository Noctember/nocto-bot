package ksoft

import (
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"strconv"
)

func Meme(ctx *gocto.CommandContext) {
	meme, err := utils.Ksoft.RandomMeme()
	if err != nil {
		ctx.Error(err)
	}
	embed := gocto.NewEmbed().
		SetTitle(meme.Title).
		SetFooter(meme.Subreddit).
		AddInlineField("Upvotes", strconv.FormatInt(int64(meme.Upvotes), 10)).
		AddInlineField("Downvotes", strconv.FormatInt(int64(meme.Downvotes), 10)).
		SetImage(meme.ImageURL)
	ctx.ReplyEmbed(embed.Build())
}
