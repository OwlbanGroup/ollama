package cmd

import (
	"regexp"
	"strings"
)

// extractFileNames extracts file names from the input string.
func extractFileNames(input string) []string {
	// Regular expression to match file paths
	re := regexp.MustCompile(`(?i)(?:[a-zA-Z]:)?(?:[\\/][^\\/]+)+`)
	matches := re.FindAllString(input, -1)

	// Clean up the matches and return
	var fileNames []string
	for _, match := range matches {
		fileNames = append(fileNames, strings.TrimSpace(match))
	}
	return fileNames
}
