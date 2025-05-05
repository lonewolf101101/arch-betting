package common

import (
	"fmt"
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
	case "application/pdf":
		return ".pdf"
	default:
		return "invalidType"
	}
}

// Source: https://gist.github.com/chadleeshaw/5420caa98498c46a84ce94cd9655287a
func FormatFileSize(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
