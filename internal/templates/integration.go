package templates

import (
	"fmt"
	"github.com/wakflo/go-sdk/integration"
	"os"
	"path/filepath"
	"strings"
)

type CreateIntegrationProps struct {
	integration.IntegrationSchemaModel
	Docs string `json:"name" toml:"name" yaml:"name"`
}

func CreateIntegrationFolder(meta *CreateIntegrationProps) error {
	// Generate folder structure
	folderName := strings.ToLower(meta.Name)
	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create folder '%s': %w", folderName, err)
	}

	// Populate the folder with boilerplate files
	files := map[string]string{
		"lib.go":        libGoTemplate,
		"README.md":     readmeTemplate,
		integrationFile: integrationTomlTemplate,
	}

	for fileName, content := range files {
		filePath := filepath.Join(folderName, fileName)
		if err := WriteTemplateToFile(filePath, content, meta); err != nil {
			return fmt.Errorf("failed to create file '%s': %w", filePath, err)
		}
	}

	fmt.Printf("Integration '%s' created successfully in folder '%s'.\n", meta.Name, folderName)
	return nil
}

// Templates for the integration files

const libGoTemplate = `package {{ .Name | toLower }}

import (
	"github.com/wakflo/go-sdk/integration"
)

var Integration = integration.Register(New{{ .Name }}())

type {{ .Name }} struct{}

func (n *{{ .Name }}) Auth() *integration.Auth {
	return &integration.Auth{
		Required: false,
	}
}

func (n *{{ .Name }}) Triggers() []integration.Trigger {
	return []integration.Trigger{}
}

func (n *{{ .Name }}) Actions() []integration.Action {
	return []integration.Action{}
}

func New{{ .Name }}() integration.Integration {
	return &{{ .Name }}{}
}
`

const readmeTemplate = `# {{ .Name }} Integration

## Description

{{ .Description }}

{{ .Docs }}

## Categories

{{ range .Categories }}- {{ . }}
{{ end }}

## Authors

{{ range .Authors }}- {{ . }}
{{ end }}
`

const integrationTomlTemplate = `[integration]
name = "{{ .Name }}"
description = "{{ .Description }}"
version = "{{ .Version }}"
icon = "{{ .Icon }}"
categories = [{{ range $index, $element := .Categories }}{{ if $index }}, {{ end }}"{{ $element }}"{{ end }}]
authors = [{{ range $index, $element := .Authors }}{{ if $index }}, {{ end }}"{{ $element }}"{{ end }}]
`
