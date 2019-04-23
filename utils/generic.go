package utils

import (
	"os"
	"strings"
)

func Capitalize(s string) string {
	var builder strings.Builder

	if s[0] > 96 && s[0] < 123 {
		builder.Grow(len(s))
		builder.WriteByte(s[0] - 32)
		builder.WriteString(s[1:])
		return builder.String()
	}

	return s
}

func IsDebug() bool {
	_, ok := os.LookupEnv("FANGS_DEV")
	return ok
}
