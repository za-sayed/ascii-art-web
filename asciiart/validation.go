package asciiart

func Validation(str string) bool {
	for _, char := range str {
		if char > 126 {
			return false
		}
	}
	return true
}
