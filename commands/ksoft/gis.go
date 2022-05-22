package ksoft

import (
	"Noctobot/utils"
	"fmt"
	ksoftgo "github.com/KSoft-Si/KSoftgo"
	"github.com/Noctember/gocto"
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

func GIS(ctx *gocto.CommandContext) {
	query := ctx.JoinedArgs()

	zoom := 12
	if ctx.HasFlag("zoom") {
		zoom = cast.ToInt(ctx.Flag("zoom"))
	}

	gis, err := utils.Ksoft.GetGIS(ksoftgo.ParamGIS{Location: query, IncludeMap: true, Zoom: zoom})
	if err != nil {
		ctx.Error(err)
		return
	}
	if gis.Error {
		ctx.Reply("%s Location not found", utils.CrossMark)
		return
	}

	embed := gocto.NewEmbed().
		SetTitle(gis.Data.Address).
		SetImage(gis.Data.Map).
		AddField("Latitude", strconv.FormatFloat(gis.Data.Lat, 'f', -1, 64)).
		AddField("Longitude", strconv.FormatFloat(gis.Data.Lon, 'f', -1, 64)).
		AddField("Tags", strings.Join(gis.Data.Type, ", ")).
		SetFooter(fmt.Sprintf("Zoom: %d", zoom))

	ctx.ReplyEmbed(embed.Build())
}
