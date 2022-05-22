package util

import (
	"Noctobot/utils"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	qr "github.com/kiktomo/goqr"
	"github.com/liyue201/goqr"
	"github.com/nfnt/resize"
	"github.com/spf13/cast"
	"image"
	"image/png"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	messagePattern = regexp.MustCompile(`(?m)https?:\/\/(?:(?:ptb|canary|development)\.)?discordapp\.com\/channels\/(?P<guild>\d{15,21})\/(?P<channel>\d{15,21})\/(?P<message>\d{15,21})\/?`)
	urlPattern     = regexp.MustCompile(`(?m)(?:https:\/\/|http:\/\/)[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&\/\/=]*(?:\.png|\.jpg|\.jpeg|\.gif|\.gifv|\.webp))`)
)

func QRCode(ctx *gocto.CommandContext) {
	bot := ctx.Bot
	ttype := ctx.Arg(0).AsString()
	content := strings.Join(ctx.RawArgs[1:], " ")
	switch ttype {
	case "decode":
		channel := ctx.Channel.ID
		message := content
		guild := ctx.Guild.ID
		if messagePattern.MatchString(content) {
			n1 := messagePattern.SubexpNames()
			r2 := messagePattern.FindAllStringSubmatch(content, -1)[0]

			md := make(map[string]interface{})
			for i, n := range r2 {
				md[n1[i]] = n
			}

			guild = cast.ToInt64(md["guild"])
			channel = cast.ToInt64(md["channel"])
			message = cast.ToString(md["message"])
		}
		if len(ctx.Message.Attachments) > 0 {
			message = string(ctx.Message.ID)
			channel = ctx.Channel.ID
			guild = ctx.Guild.ID
		}

		if !urlPattern.MatchString(content) {
			msg, _ := bot.Session.ChannelMessage(channel, cast.ToInt64(message))
			if len(msg.Attachments) < 1 {
				ctx.Reply("%s The provided message does not contain an attachment.", utils.CrossMark)
				return
			}
			for _, att := range msg.Attachments {
				res, _ := http.Get(att.URL)

				defer res.Body.Close()
				img, _, err := image.Decode(res.Body)
				if err != nil {
					ctx.Error("image.Decode error: %v\n", err)
					return
				}

				qrCodes, err := goqr.Recognize(img)
				if err != nil || len(qrCodes) < 1 {
					ctx.Reply("%s No qr code found.", utils.CrossMark)
					return
				}

				embed := gocto.NewEmbed().
					SetTitle("QRCode scanner").
					SetThumbnail(att.URL)

				for i, qrCode := range qrCodes {
					embed.AddField("Content #"+strconv.FormatInt(int64(i+1), 10), fmt.Sprintf("%s", qrCode.Payload))
				}

				embed.AddField("Where?", fmt.Sprintf("[Jump to message](https://discordapp.com/channels/%s/%s/%s)", guild, channel, message))
				bot.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed.Build())
			}
		} else {
			res, _ := http.Get(content)

			defer res.Body.Close()
			img, _, err := image.Decode(res.Body)
			if err != nil {
				ctx.Error("image.Decode error: %v\n", err)
				return
			}

			qrCodes, err := goqr.Recognize(img)
			if err != nil || len(qrCodes) < 1 {
				ctx.Reply("%s No qr code found.", utils.CrossMark)
				return
			}

			embed := gocto.NewEmbed().
				SetTitle("QRCode scanner").
				SetThumbnail(content)

			for i, qrCode := range qrCodes {
				embed.AddField("Content #"+strconv.FormatInt(int64(i+1), 10), fmt.Sprintf("%s", qrCode.Payload))
			}

			bot.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed.Build())
		}

	default:
		content = strings.Join(ctx.RawArgs, " ")
		fallthrough
	case "encode":
		code, err := qr.Encode(content, 0, qr.ECLevelM)
		f, err := os.Create("./qrcode.png")

		newImage := resize.Resize(160, 0, code, resize.NearestNeighbor)
		defer f.Close()

		png.Encode(f, newImage)
		f, _ = os.Open("qrcode.png")
		_, err = ctx.Session.ChannelMessageSendComplex(ctx.Channel.ID, &discordgo.MessageSend{
			Content: "",
			Embed: &discordgo.MessageEmbed{
				Title: "**QR Code**",
				Image: &discordgo.MessageEmbedImage{
					URL:    "attachment://qrcode.png",
					Width:  1024,
					Height: 1024,
				},
				Footer: &discordgo.MessageEmbedFooter{
					Text:    fmt.Sprintf("%s#%s", ctx.Author.Username, ctx.Author.Discriminator),
					IconURL: ctx.Author.AvatarURL("1024"),
				},
			},
			Files: []*discordgo.File{{Name: "qrcode.png", Reader: f}},
		})
		if err != nil {
			ctx.Error(err)
		}
	}
}
