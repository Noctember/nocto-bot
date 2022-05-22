package ksoft

import (
	"Noctobot/utils"
	"bufio"
	"fmt"
	ksoftgo "github.com/KSoft-Si/KSoftgo"
	"github.com/Noctember/gocto"
	"github.com/jonas747/discordgo"
	"strings"
)

func Lyrics(ctx *gocto.CommandContext) {
	query := ctx.JoinedArgs()
	if len(ctx.RawArgs) < 1 {
		pres, err := ctx.Session.State.Presence(ctx.Guild.ID, ctx.Author.ID)
		if err != nil {
			ctx.Reply("%s You aren't listening to Spotify!", utils.CrossMark)
			return
		}
		if pres.Game.Type == discordgo.GameTypeListening {
			query = pres.Game.State + " " + pres.Game.Details
		} else {
			for _, act := range pres.Activities {
				if act.Type == discordgo.GameTypeListening {
					query = act.State + " " + act.Details
				}
			}
			if len(query) < 2 || query == "" {
				ctx.Reply("%s You aren't listening to Spotify!", utils.CrossMark)
				return
			}
		}

	}
	lyrics, err := utils.Ksoft.SearchLyrics(ksoftgo.ParamSearchLyrics{Query: query})
	if err != nil {
		ctx.Error(err)
		return
	}
	if len(lyrics.Data) < 1 {
		ctx.Reply("%s No lyrics found", utils.CrossMark)
		return
	}
	p := gocto.NewPaginatorForContext(ctx)
	p.Delete()
	p.SetTemplate(func() *gocto.Embed {
		return gocto.NewEmbed().SetAuthor(fmt.Sprintf("%s by %s", lyrics.Data[0].Name, lyrics.Data[0].Artist), "https://cdn.ksoft.si/images/Logo1024%20-%20W.png", lyrics.Data[0].URL).SetThumbnail(lyrics.Data[0].AlbumArt)
	})
	p.SetExtra("| Powered by KSoft.Si")
	var out string

	scanner := bufio.NewScanner(strings.NewReader(lyrics.Data[0].Lyrics))
	for scanner.Scan() {
		if len(out) > 1000 {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.SetDescription(out)
			})
			out = ""
		}
		out += scanner.Text() + "\n"
	}
	if len(out) != 0 {
		p.AddPage(func(em *gocto.Embed) *gocto.Embed {
			return em.SetDescription(out)
		})
	}
	go p.Run()
}
