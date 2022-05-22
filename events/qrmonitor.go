package events

import (
	"Noctobot/database"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"github.com/liyue201/goqr"
	"image"
	"net/http"
)

type Data struct {
	Guild     int    `json:"guild"`
	Channel   int    `json:"channel"`
	Message   string `json:"message"`
	MessageID int    `json:"message_id"`
	User      int    `json:"user"`
}

func QRCode(bot *gocto.Bot, ctx *gocto.MonitorContext) {
	if len(ctx.Message.Attachments) < 1 {
		return
	}
	conf, err := database.GetGConfig(ctx.Guild.ID)
	if err == nil && conf.QRCode {
		for _, att := range ctx.Message.Attachments {
			res, _ := http.Get(att.URL)

			defer res.Body.Close()
			img, _, err := image.Decode(res.Body)
			if err != nil {
				fmt.Printf("image.Decode error: %v\n", err)
				return
			}
			_, err = goqr.Recognize(img)
			if err != nil {
				return
			}
			err = ctx.Session.MessageReactionAdd(ctx.Channel.ID, ctx.Message.ID, ":qrcode:678409598306222121")
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			return
		}
	}
	return
}
func Listen(bot *gocto.Bot) {
	var r *discordgo.MessageReaction

	go func() {
		for {
			select {
			case e := <-nextReaction(bot):
				r = e.MessageReaction
			}

			msg, _ := bot.Session.ChannelMessage(r.ChannelID, r.MessageID)
			author, _ := bot.Session.User(r.UserID)
			if len(msg.Attachments) < 1 || author.Bot {
				return
			}
			go func() {
				switch r.Emoji.Name {
				case "qrcode":
					channel, _ := bot.Session.Channel(r.ChannelID)
					guild, _ := bot.Session.Guild(r.GuildID)
					cmd := bot.GetCommand("qr")
					args := append([]string{}, "decode", string(msg.ID))

					cctx := &gocto.CommandContext{
						Bot:     bot,
						Message: msg,
						Command: cmd,
						Channel: channel,
						Session: bot.Session,
						Author:  author,
						RawArgs: args,
						Prefix:  "",
						Guild:   guild,
					}
					if !cctx.ParseArgs() {
						return
					}
					cmd.Run(cctx)
				}
			}()
		}
	}()
}

func nextReaction(bot *gocto.Bot) chan *discordgo.MessageReactionAdd {
	channel := make(chan *discordgo.MessageReactionAdd)
	bot.Session.AddHandlerOnce(func(_ *discordgo.Session, r *discordgo.MessageReactionAdd) {
		channel <- r
	})
	return channel
}
