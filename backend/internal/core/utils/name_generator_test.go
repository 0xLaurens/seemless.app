package utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"unicode"
)

func TestShouldGenerateRandomNameShouldReturnCamelCase(t *testing.T) {
	amount := 1000
	pattern := regexp.MustCompile("[A-Z][^A-Z]*")

	for i := 0; i < amount; i++ {
		name := GenerateRandomDisplayName()
		parts := pattern.FindAllString(name, -1)
		assert.Equal(t, 2, len(parts))

		adjective := parts[0]
		assert.Equal(t, true, unicode.IsUpper(rune(adjective[0])))
		noun := parts[1]
		assert.Equal(t, true, unicode.IsUpper(rune(noun[0])))
	}
}
