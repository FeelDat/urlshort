// Package utils provides utility functions for various common tasks.
package utils

import "strings"

// alphabet contains the set of characters used for base62 encoding.
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Base62Encode encodes the provided number into a base62 format string. The encoding
// uses a mixture of lowercase letters, uppercase letters, and numbers.
//
// Parameters:
//   - number: The unsigned 64-bit integer to encode.
//
// Returns:
//   - A base62 encoded string representing the provided number.
func Base62Encode(number uint64) string {
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}
	return encodedBuilder.String()
}
