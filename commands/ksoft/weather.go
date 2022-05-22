package ksoft

import (
	"Noctobot/utils"
	"fmt"
	ksoftgo "github.com/KSoft-Si/KSoftgo"
	"github.com/Noctember/gocto"
	"strconv"
)

func Weather(ctx *gocto.CommandContext) {
	query := ctx.JoinedArgs()
	weather, err := utils.Ksoft.GetWeather(ksoftgo.ParamWeather{Location: query, ReportType: "currently"})

	if weather.Error || err != nil {
		ctx.Reply("%s Could not find location.", utils.CrossMark)
		return
	}

	embed := gocto.NewEmbed().
		SetAuthor(weather.Data.Location.Address).
		SetThumbnail(weather.Data.IconURL).
		AddInlineField("Temperature", strconv.FormatFloat(weather.Data.Temperature, 'f', -1, 64)).
		AddInlineField("Feels like", strconv.FormatFloat(weather.Data.ApparentTemperature, 'f', -1, 64)).
		AddInlineField("Humidity", fmt.Sprintf("%.0f", weather.Data.Humidity*100)+"%").
		AddInlineField("Pressure", strconv.FormatFloat(weather.Data.Pressure, 'f', -1, 64)).
		SetDescription(weather.Data.Summary).
		SetFooter("Units: " + weather.Data.Units)

	if len(weather.Data.Alerts) > 0 {
		for _, a := range weather.Data.Alerts {
			embed.AddField(a.Title, a.Severity+" "+a.Description)
		}
	}

	_, err = ctx.ReplyEmbed(embed.Build())
}
