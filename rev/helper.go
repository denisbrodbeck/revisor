package rev

import (
	"io"
	"log"
	"strings"
)

// logCloser is a convenience function to log the error on a deferred close.
func logCloser(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("Failed to close handle: %s", err)
	}
}

func hasSuffixAny(s string, suffix ...string) bool {
	s = strings.ToLower(s)
	for _, suf := range suffix {
		if strings.HasSuffix(s, strings.ToLower(suf)) {
			return true
		}
	}
	return false
}

func filterFilesBySuffix(files Files, suffixes ...string) Files {
	res := Files{}
	for _, file := range files {
		if hasSuffixAny(file.Name, suffixes...) {
			res = append(res, file)
		}
	}
	return res
}
