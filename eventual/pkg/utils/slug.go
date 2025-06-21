package utils

import (
	"regexp"
	"strings"
)

// Slugify converts a string to a URL-friendly slug
func Slugify(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	// Replace spaces with hyphens
	text = strings.ReplaceAll(text, " ", "-")
	// Remove special characters except hyphens
	reg := regexp.MustCompile("[^a-z0-9-]")
	text = reg.ReplaceAllString(text, "")
	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile("-+")
	text = reg.ReplaceAllString(text, "-")
	// Remove leading and trailing hyphens
	text = strings.Trim(text, "-")
	return text
}
