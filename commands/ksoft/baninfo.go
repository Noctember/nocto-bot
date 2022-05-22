package ksoft

import (
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
)

func BanInfo(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsString()
	info, err := utils.Ksoft.GetBanInfo(cast.ToInt64(target))
	if !info.Exists {
		ctx.Reply("%s This user is not banned from KSoft.Si", utils.CrossMark)
		return
	}

	if err != nil {
		ctx.Error(err)
		return
	}

	embed := gocto.NewEmbed().
		SetColor(((26&0x0ff)<<16)|((201&0x0ff)<<8)|(6&0x0ff)).
		AddField("Proof", "[Link](%s)", info.Proof).
		SetTitle(info.Name+"#"+info.Discriminator).
		AddField("ID", info.ID).
		AddField("Added by", "<@%s>", info.ModeratorID).
		AddField("Reason", info.Reason).
		SetFooter("Powered by KSoft.Si")

	if info.IsBanActive {
		embed.SetColor(((201 & 0x0ff) << 16) | ((29 & 0x0ff) << 8) | (6 & 0x0ff))
	}

	ctx.BuildEmbed(embed)
}
