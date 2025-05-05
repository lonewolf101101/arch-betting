package common

import (
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var numberPrinter = message.NewPrinter(language.Mongolian)

func FormatAmount(f float32) string {
	r := numberPrinter.Sprintf("%.2f", f)
	r = strings.ReplaceAll(r, ",", " ")
	r = strings.ReplaceAll(r, ".", ",")
	return r
}

func FileFormat(contentType *string) string {
	switch *contentType {
	case "image/png":
		return ".png"
	case "image/jpg":
		return ".jpg"
	case "image/jpeg":
		return ".jpg"
	default:
		return "invalidType"
	}
}
