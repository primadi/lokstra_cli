// File: cmd/lokstra/internal/lint/service_uri_checker.go

package lint

import (
	"fmt"
	"os"
	"regexp"

	"github.com/primadi/lokstra_cli/internal/uri"
)

var serviceURIPattern = regexp.MustCompile(`lokstra://[^\s"']+`)

func CheckServiceURIFormat(filePath string, content []byte) []string {
	var issues []string
	matches := serviceURIPattern.FindAllString(string(content), -1)
	for _, match := range matches {
		if err := uri.ValidateServiceURI(match); err != nil {
			issues = append(issues, fmt.Sprintf("❌ %s: %s", filePath, err.Error()))
		}
	}
	return issues
}

func LintGoFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("❌ Failed to read %s: %v", path, err)}
	}
	return CheckServiceURIFormat(path, data)
}
