package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"io/fs"

	"github.com/primadi/lokstra_cli/internal/lint"
	"github.com/spf13/cobra"
)

const goVersion = "1.24"

// TemplateContext holds variables used in .tpl files.
type TemplateContext struct {
	AppName    string
	ModuleName string
	// Add more fields if needed in templates
}

var initCmd = &cobra.Command{
	Use:   "init [project-type] [name]",
	Short: "Initialize a new Lokstra project",
	Long: `Initialize a new Lokstra project of type:
server, module, service, middleware, plugin

Template Resolution:
The --template flag supports flexible template resolution:
1. If no template specified, uses LOKSTRA_TEMPLATE env var or "default"
2. If template is a valid directory path, uses it directly
3. Otherwise, looks for template under ./scaffold/ directory

Output Directory:
By default, creates project in ./<name>. Use --output to specify custom location.

Examples:
  lokstra init server my-app                          # Creates ./my-app/
  lokstra init server my-app --template custom        # Uses ./scaffold/custom/
  lokstra init server my-app -o /path/to/projects     # Creates /path/to/projects/my-app/
  lokstra init server my-app -o ../projects --template /custom/tpl  # Multiple flags`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		validTypes := map[string]bool{
			"server":     true,
			"module":     true,
			"service":    true,
			"middleware": true,
			"plugin":     true,
		}

		typeName := args[0]
		name := args[1]

		if !validTypes[typeName] {
			fmt.Printf("❌ Invalid project type: %s\n", typeName)
			fmt.Println("Valid types are: server, module, service, middleware, plugin")
			os.Exit(1)
		}

		module, _ := cmd.Flags().GetString("module")
		if module == "" {
			module = "github.com/example/" + name
		}

		template, _ := cmd.Flags().GetString("template")
		output, _ := cmd.Flags().GetString("output")

		if err := initGenericProject(typeName, name, module, template, output); err != nil {
			log.Fatalf("❌ Failed to initialize %s: %v", typeName, err)
		}

		outputPath := output
		if outputPath == "" {
			outputPath = "."
		}
		projectPath := filepath.Join(outputPath, name)

		fmt.Printf("✅ Lokstra %s project created: %s\n", typeName, projectPath)
		if typeName == "server" {
			fmt.Printf("🚀 Try running: cd %s && go run cmd/main.go\n", projectPath)
		}
	},
}

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint the current Lokstra project",
	Run: func(cmd *cobra.Command, args []string) {
		var goFiles, yamlFiles []string
		excludedDirs := map[string]bool{
			".git": true, "vendor": true, "node_modules": true, "bin": true, "dist": true,
		}

		err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() && excludedDirs[d.Name()] {
				return filepath.SkipDir
			}
			if strings.HasSuffix(path, ".go") {
				goFiles = append(goFiles, path)
			} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
				yamlFiles = append(yamlFiles, path)
			}
			return nil
		})
		if err != nil {
			log.Fatalf("❌ Lint error: %v", err)
		}

		fmt.Printf("🔍 Found %d Go files, %d YAML files\n", len(goFiles), len(yamlFiles))

		totalIssues := 0

		for _, file := range goFiles {
			issues := lint.LintGoFile(file)
			for _, issue := range issues {
				fmt.Println(issue)
			}
			totalIssues += len(issues)
		}

		for _, file := range yamlFiles {
			issues := lint.LintYamlFile(file)
			for _, issue := range issues {
				fmt.Println(issue)
			}
			totalIssues += len(issues)
		}

		if totalIssues == 0 {
			fmt.Println("✅ No lint issues found")
		} else {
			fmt.Printf("❌ Total lint issues: %d\n", totalIssues)
			os.Exit(1)
		}
	},
}

func init() {
	initCmd.Flags().String("module", "", "Go module name (default: github.com/example/<name>)")
	initCmd.Flags().String("template", "", "Template directory path or name. Searches: 1) Direct path if valid directory, 2) ./scaffold/<template>/, 3) LOKSTRA_TEMPLATE env var, 4) 'default'")
	initCmd.Flags().StringP("output", "o", "", "Output directory for the project (default: current directory)")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(lintCmd)
}

// initGenericProject handles any project type by using the matching scaffold directory.
func initGenericProject(typeName, name, module, template, output string) error {
	ctx := TemplateContext{
		AppName:    name,
		ModuleName: module,
	}

	// Determine project root directory
	var root string
	if output == "" {
		root = filepath.Join(".", name)
	} else {
		root = filepath.Join(output, name)
	}

	// Create the project directory
	if err := os.MkdirAll(root, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	templatePath, err := resolveTemplatePath(template)
	if err != nil {
		return fmt.Errorf("failed to resolve template path: %w", err)
	}

	templateRoot := filepath.Join(templatePath, typeName)

	gomod := fmt.Sprintf("module %s\n\ngo %s\n", module, goVersion)
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(gomod), 0644); err != nil {
		return err
	}

	err = filepath.Walk(templateRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(templateRoot, path)
		if err != nil {
			return err
		}

		outPath := filepath.Join(root, relPath)
		return renderTemplateFile(path, outPath, ctx)
	})
	if err != nil {
		return err
	}

	cmdGet := exec.Command("go", "get", "github.com/primadi/lokstra@latest")
	cmdGet.Dir = root
	cmdGet.Stdout = os.Stdout
	cmdGet.Stderr = os.Stderr
	if err := cmdGet.Run(); err != nil {
		return fmt.Errorf("failed to run 'go get': %w", err)
	}

	cmdTidy := exec.Command("go", "mod", "tidy")
	cmdTidy.Dir = root
	cmdTidy.Stdout = os.Stdout
	cmdTidy.Stderr = os.Stderr
	if err := cmdTidy.Run(); err != nil {
		return fmt.Errorf("failed to run 'go mod tidy': %w", err)
	}

	return nil
}

// renderTemplateFile renders a template or copies a raw file depending on extension.
func renderTemplateFile(tplPath, outPath string, ctx any) error {
	if strings.HasSuffix(tplPath, ".tpl") {
		tplBytes, err := os.ReadFile(tplPath)
		if err != nil {
			return err
		}
		tpl, err := template.New(filepath.Base(tplPath)).Parse(string(tplBytes))
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(outPath[:len(outPath)-4]), 0755); err != nil {
			return err
		}
		outFile, err := os.Create(outPath[:len(outPath)-4])
		if err != nil {
			return err
		}
		defer outFile.Close()
		return tpl.Execute(outFile, ctx)
	}

	return copyFile(tplPath, outPath)
}

// copyFile copies a static file from src to dst.
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func resolveTemplatePath(template string) (string, error) {
	// Template resolution strategy:
	// 1. If template is empty, check LOKSTRA_TEMPLATE environment variable, otherwise use "default"
	if template == "" {
		if envPath := os.Getenv("LOKSTRA_TEMPLATE"); envPath != "" {
			template = envPath
		} else {
			template = "default"
		}
	}

	// 2. Check if template is a valid directory path (absolute or relative)
	if stat, err := os.Stat(template); err == nil && stat.IsDir() {
		return filepath.Abs(template)
	}

	// 3. Fallback to built-in templates under ./scaffold/ directory
	builtin := filepath.Join("scaffold", template)
	if stat, err := os.Stat(builtin); err == nil && stat.IsDir() {
		return builtin, nil
	}

	return "", fmt.Errorf("template '%s' not found (searched: direct path and ./scaffold/%s/)", template, template)
}
