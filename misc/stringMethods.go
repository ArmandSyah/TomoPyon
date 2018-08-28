package misc

import (
	"fmt"
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

func ReplaceSubstr(input, expr string) (output string) {
	r := regexp.MustCompile(expr)
	output = r.ReplaceAllString(input, "${1}")
	fmt.Println(output)
	return
}

func StripWhitespace(input string) string {
	return strings.Replace(input, " ", "", -1)
}
