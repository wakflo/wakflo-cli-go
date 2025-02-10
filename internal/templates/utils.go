package templates

import (
	"github.com/samber/lo"
	"os"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"toLower":     strings.ToLower,
	"toUpper":     strings.ToUpper,
	"toTitle":     strings.ToTitle,
	"toSnakeCase": lo.SnakeCase,
	"toKebabCase": lo.KebabCase,
	"toCamelCase": lo.CamelCase,
	"toPascal":    lo.PascalCase,
}

func WriteTemplateToFile(filePath, tmpl string, meta any) error {
	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse and execute template
	t := template.Must(template.New(filePath).Funcs(funcMap).Parse(tmpl))
	if err := t.Execute(file, meta); err != nil {
		return err
	}

	return nil
}

// IsIntegrationProject Check whether the current directory is an integration project folder
func IsIntegrationProject() bool {
	// Example: Check if a specific configuration file exists
	// Replace "integration.config" with the appropriate file or folder to check
	_, err := os.Stat("integration.toml")
	return err == nil
}
