package lint

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func CheckYamlSyntax(filePath string, content []byte) []string {
	var issues []string
	var out map[string]any
	if err := yaml.Unmarshal(content, &out); err != nil {
		issues = append(issues, fmt.Sprintf("❌ Invalid YAML syntax in %s: %v", filePath, err))
	}
	return issues
}

func LintYamlFile(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return []string{fmt.Sprintf("❌ Failed to read YAML %s: %v", path, err)}
	}

	return CheckYamlSyntax(path, data)
}
