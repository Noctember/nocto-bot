package info

import (
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/dustin/go-humanize"
	"sort"
)

func User(ctx *gocto.CommandContext) {
	target := ctx.Author

	if ctx.Arg(0).IsProvided() {
		target = ctx.Arg(0).AsUser()
	}
	created, _ := utils.SnowflakeTimestamp(target.ID)
	embed := gocto.NewEmbed().
		SetTitle(fmt.Sprintf("%s#%s (%d)", target.Username, target.Discriminator, target.ID)).
		SetThumbnail(target.AvatarURL("512")).
		AddField("Created", fmt.Sprintf("%s (%s)", created.Format("2 January 2006"), humanize.Time(created)))

	mem, err := ctx.Session.GuildMember(ctx.Guild.ID, target.ID)
	if err == nil {
		mems := ctx.Guild.Members
		sort.Slice(mems, func(i, j int) bool {
			return mems[i].JoinedAt < mems[j].JoinedAt
		})

		parsed, err := mem.JoinedAt.Parse()
		if err == nil {
			embed.AddField("Join date", fmt.Sprintf("%s (%s)", parsed.Format("2 January 2006"), humanize.Time(parsed))).
				AddField("Join position", fmt.Sprintf("%d", SliceIndex(len(ctx.Guild.Members), func(i int) bool { return ctx.Guild.Members[i].User.ID == target.ID })+1))
		}
	}
	pres, err := ctx.Session.State.Presence(ctx.Guild.ID, target.ID)
	if err == nil {
		embed.AddField("Status", string(pres.Status))
		if pres.Game != nil && pres.Game.Name != "" && pres.Game.State != "" {
			embed.AddField(pres.Game.Name, pres.Game.State)
		}
	}

	_, err = ctx.ReplyEmbed(embed.Build())
	if err != nil {
		ctx.Error(err)
	}

}

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}
