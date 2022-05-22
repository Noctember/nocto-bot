package hypixel

import (
	"Noctobot/utils"
	"Noctobot/utils/mcfont"
	"encoding/json"
	"fmt"
	"github.com/Noctember/gocto"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func Hypixel(ctx *gocto.CommandContext) {
	target := ctx.Arg(0).AsString()
	player := getPlayer(target)
	if !player.Success {
		ctx.Reply("%s `%s` is not a valid player.", utils.CrossMark, target)
		return
	}

	msg, _ := mcfont.Parse(player.Player.Display)
	var formatted string
	for _, f := range msg {
		formatted += f.Message
	}
	//created.Format("2 January 2006"), humanize.Time(created))
	last := time.Unix(player.Player.LastLogin/1000, 0)
	first := time.Unix(player.Player.FirstLogin/1000, 0)
	if player.Player.Version == "" {
		player.Player.Version = "Unknown"
	}

	status := "Offline!"
	color := ((201 & 0x0ff) << 16) | ((29 & 0x0ff) << 8) | (6 & 0x0ff)
	if player.Player.LastLogout != 0 && player.Player.LastLogin > player.Player.LastLogout {
		status = "Online!"
		color = ((26 & 0x0ff) << 16) | ((201 & 0x0ff) << 8) | (6 & 0x0ff)
	}

	embed := gocto.NewEmbed().
		SetThumbnail(fmt.Sprintf("https://visage.surgeplay.com/full/256/%s", player.Player.UUID)).
		SetTitle(formatted).
		SetColor(color).
		AddInlineField("Status", status).
		AddInlineField("Network Level", "`%s`", strconv.FormatFloat(player.Player.NetworkLevel, 'f', 0, 64)).
		AddInlineField("Karma", "`%s`", strconv.FormatInt(int64(player.Player.Karma), 10)).
		AddInlineField("MC Version", "`%s`", player.Player.Version).
		AddInlineField("Last login", "`%s (%s)`", last.Format("2 January 2006"), humanize.Time(last)).
		AddInlineField("First login", "`%s (%s)`", first.Format("2 January 2006"), humanize.Time(first)).
		SetFooter(fmt.Sprintf("%s's stats on Hypixel.net", player.Player.DisplayName))

	ctx.ReplyEmbed(embed.Build())
}

func getPlayer(name string) Player {
	res, err := http.Get("https://api.sk1er.club/player/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resp := Player{}
	by, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(by, &resp)
	return resp
}
