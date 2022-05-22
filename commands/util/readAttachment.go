package util

import (
	"Noctobot/utils"
	"bufio"
	"bytes"
	"github.com/Noctember/gocto"
	"github.com/gabriel-vasile/mimetype"
	"github.com/jonas747/discordgo"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var sitePattern = regexp.MustCompile(`(?m)^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)

func ReadAttachment(ctx *gocto.CommandContext) {
	target := ctx.JoinedArgs()

	channel := ctx.Channel.ID
	message := target

	if messagePattern.MatchString(target) {
		n1 := messagePattern.SubexpNames()
		r2 := messagePattern.FindAllStringSubmatch(target, -1)[0]

		md := make(map[string]interface{})
		for i, n := range r2 {
			md[n1[i]] = n
		}

		channel = cast.ToInt64(md["channel"])
		message = cast.ToString(md["message"])
	}

	p := gocto.NewPaginatorForContext(ctx)
	p.Delete()
	p.SetTemplate(func() *gocto.Embed {
		return gocto.NewEmbed().SetAuthor("Attachment viewer")
	})

	if !messagePattern.MatchString(target) {
		res, _ := http.Get(target)
		var buf bytes.Buffer
		tee := io.TeeReader(res.Body, &buf)

		mime, err := mimetype.DetectReader(tee)
		if err != nil {
			ctx.Error(err)
			return
		}
		if strings.Contains(mime.String(), "text") {
			var out string
			scanner := bufio.NewScanner(&buf)
			for scanner.Scan() {
				if len(out) > 600 {
					p.AddPage(func(em *gocto.Embed) *gocto.Embed {
						return em.AddField(":spider_web: "+target, "```\n"+out+"```")
					})
					out = ""
				}
				out += scanner.Text() + "\n"
			}

			if out != "" {
				p.AddPage(func(em *gocto.Embed) *gocto.Embed {
					return em.AddField(":spider_web: "+target, "```\n"+out+"```")
				})
			}
		} else if strings.Contains(mime.String(), "image") {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.AddField(":camera: "+target, " ឵឵").SetImage(target)
			})
		}
		res.Body.Close()
	} else {
		msg, _ := ctx.Session.ChannelMessage(channel, cast.ToInt64(message))

		if len(msg.Attachments) < 1 {
			ctx.Reply("%s This message does not have attachments", utils.CrossMark)
			return
		}

		for _, att := range msg.Attachments {
			attRead(att, p, ctx)
		}
	}

	go p.Run()
}

func attRead(att *discordgo.MessageAttachment, p *gocto.Paginator, ctx *gocto.CommandContext) {
	res, _ := http.Get(att.URL)
	var buf bytes.Buffer
	tee := io.TeeReader(res.Body, &buf)

	mime, err := mimetype.DetectReader(tee)
	if err != nil {
		ctx.Error(err)
		return
	}
	if strings.Contains(mime.String(), "text") {
		var out string
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			if len(out) > 600 {
				p.AddPage(func(em *gocto.Embed) *gocto.Embed {
					return em.AddField(":paperclip: "+att.Filename, "```\n"+out+"```")
				})
				out = ""
			}
			out += scanner.Text() + "\n"
		}

		if out != "" {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.AddField(":paperclip: "+att.Filename, "```\n"+out+"```")
			})
		}
	} else if strings.Contains(mime.String(), "image") {
		p.AddPage(func(em *gocto.Embed) *gocto.Embed {
			return em.AddField(":camera: "+att.Filename, " ឵឵").SetImage(att.URL)
		})
	}
	res.Body.Close()
}
