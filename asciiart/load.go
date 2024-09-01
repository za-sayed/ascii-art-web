package asciiart

import (
	"bufio"
	"os"
)

func Load(filePath string) (map[rune][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	asciiArt := make(map[rune][]string)
	var lines []string
	char := rune(32)
	for i, line := range fileLines {
		if i % 9 == 0 {
			continue
		}
		if len(lines) < 9 {
			lines = append(lines, line)
		} 
		if len(lines) == 8 {
			asciiArt[char] = lines
			lines = nil
			char++
		}
	}
	return asciiArt, nil
}

