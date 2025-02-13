package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wakflo/go-sdk/sdk"
)

type CreateIntegrationProps struct {
	sdk.IntegrationSchemaModel
	Docs string `json:"name" toml:"name" yaml:"name"`
}

func CreateIntegrationFolder(meta *CreateIntegrationProps) error {
	// Generate folder structure
	folderName := strings.ReplaceAll(strings.ToLower(meta.Name), " ", "")
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

const libGoTemplate = `package {{ .Name | toPackageName }}

import (
	_ "embed"

	"github.com/wakflo/go-sdk/sdk"
)

//go:embed README.md
var ReadME string

//go:embed flo.toml
var Flow string

var Integration = sdk.Register(New{{ .Name | toPascal }}(), Flow, ReadME)

type {{ .Name | toPascal }} struct{}

func (n *{{ .Name | toPascal }}) Auth() *sdk.Auth {
	return &sdk.Auth{
		Required: false,
	}
}

func (n *{{ .Name | toPascal }}) Triggers() []sdk.Trigger {
	return []sdk.Trigger{}
}

func (n *{{ .Name | toPascal }}) Actions() []sdk.Action {
	return []sdk.Action{}
}

func New{{ .Name | toPascal }}() sdk.Integration {
	return &{{ .Name | toPascal }}{}
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
description = """{{ .Description }}"""
version = "{{ .Version }}"
icon = "{{ .Icon }}"
categories = [{{ range $index, $element := .Categories }}{{ if $index }}, {{ end }}"{{ $element }}"{{ end }}]
authors = [{{ range $index, $element := .Authors }}{{ if $index }}, {{ end }}"{{ $element }}"{{ end }}]
`
