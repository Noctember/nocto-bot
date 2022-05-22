package admin

import (
	"bufio"
	"fmt"
	"github.com/Noctember/gocto"
	"sort"
	"strings"
)

func Listguilds(ctx *gocto.CommandContext) {
	sorting := "members"
	if ctx.HasArgs() {
		sorting = ctx.Arg(0).AsString()
	}
	var str []string
	switch sorting {
	case "joined":
		sort.Slice(ctx.Session.State.Guilds, func(i, j int) bool {
			return ctx.Session.State.Guilds[i].JoinedAt > ctx.Session.State.Guilds[j].JoinedAt
		})
		for _, k := range ctx.Session.State.Guilds {
			time, _ := k.JoinedAt.Parse()
			str = append(str, fmt.Sprintf("%s (Joined: %s)", k.Name, time.Format("2 January 2006")))
		}
	default:
		fallthrough
	case "members":
		sort.Slice(ctx.Session.State.Guilds, func(i, j int) bool {
			return ctx.Session.State.Guilds[i].MemberCount > ctx.Session.State.Guilds[j].MemberCount
		})
		for _, k := range ctx.Session.State.Guilds {
			str = append(str, fmt.Sprintf("%s (Members: %d)", k.Name, k.MemberCount))
		}
	}
	p := gocto.NewPaginatorForContext(ctx)
	p.Delete()
	p.SetTemplate(func() *gocto.Embed {
		return gocto.NewEmbed().SetAuthor("Nocto List guilds")
	})

	var out string
	scanner := bufio.NewScanner(strings.NewReader(strings.Join(str, "\n")))
	for scanner.Scan() {
		if len(out) > 600 {
			p.AddPage(func(em *gocto.Embed) *gocto.Embed {
				return em.AddField("📤 Guilds", "```\n"+out+"```")
			})
			out = ""
		}
		out += scanner.Text() + "\n"
	}

	if out != "" {
		p.AddPage(func(em *gocto.Embed) *gocto.Embed {
			return em.AddField("📤 Guilds", "```\n"+out+"```")
		})
	}
	go p.Run()

}
