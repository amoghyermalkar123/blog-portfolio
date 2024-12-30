package utils

import (
	"regexp"
	"strings"
)

type TOCEntry struct {
	ID       string
	Level    int
	Title    string
	Children []TOCEntry
}

// CalculateReadingTime estimates reading time in minutes
func CalculateReadingTime(content string) int {
	words := len(strings.Fields(content))
	wordsPerMinute := 200                                    // average reading speed
	minutes := (words + wordsPerMinute - 1) / wordsPerMinute // round up
	if minutes < 1 {
		return 1
	}
	return minutes
}

// GenerateTableOfContents parses markdown content for headings
func GenerateTableOfContents(content string) []TOCEntry {
	var toc []TOCEntry

	// Regex for markdown headings
	headingRegex := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	matches := headingRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		level := len(match[1]) // number of # symbols
		title := strings.TrimSpace(match[2])
		id := strings.ToLower(strings.ReplaceAll(title, " ", "-"))

		// Remove any special characters from ID
		id = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(id, "")

		entry := TOCEntry{
			ID:    id,
			Level: level,
			Title: title,
		}

		toc = append(toc, entry)
	}

	return toc
}
