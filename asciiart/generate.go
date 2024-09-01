package asciiart

import (
	"strings"
)

func GenerateAsciiArt(text string, style map[rune][]string) string {
	var result []string
	parts := strings.Split(text, "\n") 
	for _, part := range parts {
		if part == "" {
			result = append(result, "")
			continue
		}
		for i := 0; i < 8; i++ {
			var lineOutput []string
			for _, char := range part {
				charLines := style[char]
				if charLines == nil {
					continue 
				}
				lineOutput = append(lineOutput, charLines[i])
			}
			result = append(result, strings.Join(lineOutput, ""))
		}
		result = append(result, "") 
	}
	return strings.Join(result, "\n")
}
