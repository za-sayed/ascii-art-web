package asciiart

import (
	"strings"
)

func GenerateAsciiArt(text string, style map[rune][]string) string {
	var result []string
	parts := splitText(text)
	for _, part := range parts {
		if part == "" {
			result = append(result, "")
			continue
		}
		result = append(result, renderPart(part, style)...)
		result = append(result, "")
	}
	return strings.Join(result, "\n")
}

func splitText(text string) []string {
	return strings.Split(text, "\n")
}

func renderPart(part string, style map[rune][]string) []string {
	var result []string
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
	return result
}
