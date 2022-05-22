package colorize

import (
	"github.com/lucasb-eyer/go-colorful"
	"regexp"
	"strconv"
	"strings"
)

const (
	rgbString                    = "rgb(%d,%d,%d)"
	rgbCaptureRegexString        = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	rgbCaptureRegexPercentString = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])%\\s*\\)$"
)

var (
	rgbCaptureRegex        = regexp.MustCompile(rgbCaptureRegexString)
	rgbCapturePercentRegex = regexp.MustCompile(rgbCaptureRegexPercentString)
)

func ParseRGB(color *colorful.Color, s string) error {

	s = strings.ToLower(s)

	vals := rgbCaptureRegex.FindAllStringSubmatch(s, -1)

	color.R, _ = strconv.ParseFloat(vals[0][1], 64)
	color.G, _ = strconv.ParseFloat(vals[0][2], 64)
	color.B, _ = strconv.ParseFloat(vals[0][3], 64)

	color.R = (color.R - 0.5) / 255.0
	color.G = (color.G - 0.5) / 255.0
	color.B = (color.B - 0.5) / 255.0

	return nil
}
