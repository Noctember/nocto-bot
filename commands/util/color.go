package util

import (
	"Noctobot/utils"
	"Noctobot/utils/colorize"
	"bytes"
	_ "encoding/hex"
	"github.com/Noctember/gocto"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/draw"
	"image/png"
	"strconv"
	"strings"
)

func Color(ctx *gocto.CommandContext) {
	co := strings.Join(ctx.RawArgs, "")
	var err error = nil
	c := colorful.Color{}
	if strings.HasPrefix(co, "#") {
		c, err = colorful.Hex(co)
	} else if strings.HasPrefix(co, "rgb") {
		println(co)
		err := colorize.ParseRGB(&c, co)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	if err != nil {
		ctx.Reply("%s Invalid Color Hex.", utils.CrossMark)
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.ZP, draw.Src)
	b := &bytes.Buffer{}
	png.Encode(b, img)
	ctx.SendFile("color.png", b, "**Hex:** %s\n**RGB:** rgb(%s, %s, %s)", c.Hex(), strconv.FormatFloat(c.R*255.0+0.5, 'f', 0, 64), strconv.FormatFloat(c.G*255.0+0.5, 'f', 0, 64), strconv.FormatFloat(c.B*255.0+0.5, 'f', 0, 64))
}
