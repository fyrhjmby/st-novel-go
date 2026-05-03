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

// countWordsFromHTML 去除 HTML 标签后统计纯文本字数
func countWordsFromHTML(html string) int {
	text := tagRegex.ReplaceAllString(html, "")
	// 按中文字符和非中文词统计
	words := 0
	for _, r := range text {
		if r > 127 { // 中文字符每个计为1字
			words++
		}
	}
	// 英文/数字按空格分词统计
	englishWords := strings.Fields(tagRegex.ReplaceAllString(html, " "))
	for _, w := range englishWords {
		isAscii := true
		for _, r := range w {
			if r > 127 {
				isAscii = false
				break
			}
		}
		if isAscii && len(w) > 0 {
			words++
		}
	}
	return words
}
