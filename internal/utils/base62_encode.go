package utils

import "strings"

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Base62Encode(number uint64) string {

	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}
	return encodedBuilder.String()
}
