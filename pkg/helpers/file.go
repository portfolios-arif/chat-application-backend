package helpers

func ValidFileType(fileType string) bool {
	validTypes := []string{"image/jpg", "image/jpeg", "image/png", "image/webp"}
	for _, t := range validTypes {
		if fileType == t {
			return true
		}
	}
	return false
}
