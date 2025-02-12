package templates

import (
	"os"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

var funcMap = template.FuncMap{
	"toLower":       strings.ToLower,
	"toUpper":       strings.ToUpper,
	"toTitle":       strings.ToTitle,
	"toSnakeCase":   lo.SnakeCase,
	"toKebabCase":   lo.KebabCase,
	"toCamelCase":   lo.CamelCase,
	"toPascal":      lo.PascalCase,
	"toPackageName": ToPackageName,
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

func ToPackageName(value string) string {
	return strings.ToLower(strings.ReplaceAll(value, " ", ""))
}
