package utils

import (
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
