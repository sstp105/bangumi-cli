package libs

import "strings"

func SplitToSlice(s string, d string) []string {
	parts := strings.Split(s, d)

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}
