package mod

import (
	"Noctobot/utils"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"math"
	"sort"
	"time"
)

func Purge(ctx *gocto.CommandContext) {
	count := ctx.Arg(0).AsInt64()
	if count > 100 {
		ctx.Reply("%s Count must be under or equal to 100", utils.CrossMark)
		return
	}

	channel := ctx.Channel
	bots := false
	images := false
	authorID := int64(0)
	for i := 0; i < len(ctx.Args); i++ {
		switch ctx.Arg(i).Value.(type) {
		case *discordgo.Channel:
			channel = ctx.Arg(i).AsChannel()
			break
		case *discordgo.User:
			authorID = ctx.Arg(i).AsUser().ID
			break
		}
	}
	if ctx.HasFlag("bots") {
		bots = true
	}
	if ctx.HasFlag("images") {
		images = true
	}

	before := ctx.Message.ID
	oldMessageIDs := make([]int, 0, count)
	newMessageIDs := make([]int, 0, count)

	for count > 0 {
		messages, err := ctx.Session.ChannelMessages(channel.ID, int(count), before, 0, 0)
		if err != nil {
			ctx.Error(err)
			return
		}

		for _, message := range messages {
			before = message.ID

			if bots && !message.Author.Bot {
				continue
			}

			if images && (len(message.Attachments) == 0 || message.Attachments[0].Height == 0) {
				continue
			}

			if authorID != 0 && authorID != message.Author.ID {
				continue
			}

			t, _ := message.Timestamp.Parse()

			if time.Since(t) > (2 * 7 * 24 * time.Hour) {
				oldMessageIDs = append(oldMessageIDs, int(message.ID))
			} else {
				newMessageIDs = append(newMessageIDs, int(message.ID))
			}
			count--

			if count == 0 {
				break
			}
		}

		if len(messages) != 100 {
			break // end of channel
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(newMessageIDs)))
	sort.Sort(sort.Reverse(sort.IntSlice(oldMessageIDs)))

	for i := 0; i < len(newMessageIDs); i += 100 {
		c := int(math.Min(float64(len(newMessageIDs)), 100))
		new := make([]int64, 0, count)
		for _, i := range newMessageIDs {
			new = append(new, int64(i))
		}
		err := ctx.Session.ChannelMessagesBulkDelete(channel.ID, new[i:i+c])
		if err != nil {
			continue
		}
	}

	_, err := ctx.Reply("☑ Successfully purged %d messages", len(newMessageIDs))
	if err != nil {
		ctx.Error(err)
		return
	}
}
