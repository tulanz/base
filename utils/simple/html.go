package simple

import (
	"regexp"
	"strings"
)

func ClearHtml(src string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "")
	return strings.TrimSpace(src)
}
