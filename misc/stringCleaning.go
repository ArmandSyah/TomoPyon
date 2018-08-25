package misc

import (
	"github.com/grokify/html-strip-tags-go"
)

func StripHTML(oldStr string) (newStr string) {
	newStr = strip.StripTags(oldStr)
	return
}
