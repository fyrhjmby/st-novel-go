package service

import (
	"regexp"
	"strings"
)

var h1Regex = regexp.MustCompile(`(?i)<h1[^>]*>(.*?)<\/h1>`)

// SyncTitleFromContent extracts the content of the first <h1> tag from HTML content.
// If an h1 tag is found, its inner text becomes the new title.
// Otherwise, the original title is returned.
func SyncTitleFromContent(htmlContent string, currentTitle string) string {
	matches := h1Regex.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		// The first submatch is the full match, the second is the content of the capture group.
		newTitle := stripHtmlTags(matches[1])
		if newTitle != "" {
			return newTitle
		}
	}
	return currentTitle
}

// A very simple function to remove any remaining HTML tags from the h1 content.
var tagRegex = regexp.MustCompile(`<[^>]*>`)

func stripHtmlTags(input string) string {
	return strings.TrimSpace(tagRegex.ReplaceAllString(input, ""))
}
