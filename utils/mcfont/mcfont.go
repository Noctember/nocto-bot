package mcfont

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"io/ioutil"
	"strings"
)

var ()

type Message struct {
	Width      int
	Height     int
	lastColour string
	Bold       bool
	Italics    bool
	Message    string
	Newx       int
	Posheight  int
	Newline    int
}

func Parse(message string) ([]Message, int) {
	slice := strings.Split(message, "")
	var result []Message
	totalWidth := 0

	bold := false
	italics := false
	lastColour := "r"
	posheight := 0
	newline := 0
	startline := true
	newx := 0
	for i, m := range slice {
		if m == "§" || m == "|" || m == "\\" {
			continue
		} else if i != 0 && slice[i-1] == "§" {
			if strings.ContainsAny(m, "01234567890abcdef") {
				lastColour = m
				bold = false
				italics = false
			}
			if m == "l" {
				bold = true
			}
			if m == "o" {
				italics = true
			}
			if m == "r" {
				bold = false
				italics = false
				lastColour = m
			}
			continue
		}
		if i != 0 && slice[i-1] == "|" || m == "n" && slice[i-1] == "\\" {
			posheight += 42
			newline += 1
			startline = true
			newx = 5
			continue
		}

		if startline {
			startline = false
		} else {
			newx += 1
		}

		face := truetype.NewFace(getFont(bold, italics), &truetype.Options{Size: 16})

		width := font.MeasureString(face, m)
		totalWidth = +width.Round()
		result = append(result, Message{width.Round(), 16, lastColour, bold, italics,
			m, newx, posheight, newline})
	}
	return result, totalWidth
}

func getFont(bold bool, italics bool) (f *truetype.Font) {
	breg, err := ioutil.ReadFile("./utils/mcfont/fonts/regular.ttf")
	if err != nil {
		panic(err)
	}
	reg, err := truetype.Parse(breg)
	if err != nil {
		panic(err)
	}
	bboldItalics, _ := ioutil.ReadFile("./utils/mcfont/fonts/bold-italics.ttf")
	boldItalics, _ := truetype.Parse(bboldItalics)
	f = reg
	if bold && italics {
		f = boldItalics
	}

	return
}
