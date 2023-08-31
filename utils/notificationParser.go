package utils

import (
	"regexp"
)

func ExtractEmailsFromNotification(notification string) []string {
	// Define a regular expression pattern for matching email addresses
	pattern := `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find all matches in the text
	return re.FindAllString(notification, -1)
}
