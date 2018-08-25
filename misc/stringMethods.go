package misc

import (
	"github.com/grokify/html-strip-tags-go"
	"regexp"
	"strings"
)

func StripHTML(oldStr string) (newStr string) {
	newStr = strip.StripTags(oldStr)
	return
}

func ExtractSubstr(input string, expr string) (output string) {
	r, _ := regexp.Compile(expr)
	output = r.FindString(input)
	return
}

func TrimSides(input, prefix, suffix string) (output string) {
	output = strings.TrimPrefix(input, prefix)
	output = strings.TrimSuffix(output, suffix)
	return
}
