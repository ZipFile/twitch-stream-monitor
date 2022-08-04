package utils

import (
	"regexp"
	"strconv"
	"strings"
)

// Make a value suitable for using as bind address for http server.
func MakeAddr(host string, port uint64) string {
	return strings.Join([]string{host, strconv.FormatUint(port, 10)}, ":")
}

func OrStr(a, b string) string {
	if a == "" {
		return b
	}
	return a
}

var ipv4Regexp = regexp.MustCompile(`\d+[-\.]\d+[-\.]\d+[-\.]\d+`)

func ObfuscateUrl(url string) string {
	if ipv4Regexp.MatchString(url) {
		return "***"
	}

	return url
}
