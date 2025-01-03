package review

import (
	"context"
	"strings"
)

type Handler func(ctx context.Context, data string) error

var (
	beginTag = `<output>`
	endTag   = `</output>`
)

func GetDefaultHandler(send func(data string)) Handler {
	var (
		cache string

		parseBegin bool
		parseEnd   bool
	)
	return func(ctx context.Context, data string) error {
		var result string
		if parseBegin && parseEnd {
			return nil
		}
		cache += data
		if !parseBegin {
			if strings.Contains(cache, beginTag) {
				parseBegin = true
				start := strings.Index(cache, beginTag) + len(beginTag)
				end := len(cache)
				if strings.Contains(cache, endTag) {
					parseEnd = true
					end = strings.Index(cache, endTag)
				}
				result = cache[start:end]
				send(result)
				cache = ""
			}
		} else {
			if strings.Contains(cache, endTag) {
				parseEnd = true
				result = cache[:strings.Index(cache, endTag)]
				cache = ""
			} else if hasOverlapEfficient(cache, endTag) {
				// has overlap, keep cache
			} else {
				send(result)
				cache = ""
			}
		}
		return nil
	}
}
func hasOverlapEfficient(a, b string) bool {
	minLen := min(len(a), len(b))
	for i := minLen; i > 0; i-- {
		if strings.HasSuffix(a, b[:i]) {
			return true
		}
	}
	return false
}
