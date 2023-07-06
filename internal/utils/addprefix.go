package utils

import "strings"

func AddPrefix(addr string) string {
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		addr = "http://" + addr
	}
	return addr
}
