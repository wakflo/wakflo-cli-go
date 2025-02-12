package templates

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/wakflo/go-sdk/client"
	"github.com/wakflo/go-sdk/sdk"
)

const integrationFile = "flo.toml"
const readmeFile = "README.md"
const libFile = "lib.go"

type ActionTriggerMetadata struct {
	Name        string
	Description string
	Type        string // ActionType or TriggerType
	TypeName    string // sdkcore.ActionType or sdkcore.TriggerType as string
	FileName    string // File-safe name
	Constructor string // Function to append (e.g. actions.NewRunPythonAction())
	Kind        string // either "action" or "trigger"
}

func HandleAddResource(kind string, cmd *cobra.Command, floClient *client.Client) error {
	// Ensure the command is being run from an integration folder
	data, err := os.ReadFile(integrationFile)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("not an integration project: missing '%s' file", integrationFile)
	}

	var schema sdk.IntegrationSchemaModel
	if err := toml.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse '%s' file: %w", integrationFile, err)
	}

	if _, err := os.Stat(libFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("missing 'lib.go' file in the integration project")
	}

	meta, err := collectInput(kind, &schema, floClient)
	if err != nil {
		return err
	}

	// Create resource folder
	resourceFolder := kind + "s"
	if _, err := os.Stat(resourceFolder); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(resourceFolder, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create resource folder: %w", err)
		}
	}

	// Create resource file
	resourceFileName := filepath.Join(resourceFolder, meta.FileName+".go")
	if err := WriteTemplateToFile(resourceFileName, getResourceTemplate(kind), meta); err != nil {
		return fmt.Errorf("failed to create %s: %w", kind, err)
	}

	// Create documentation (Markdown) file
	docFileName := filepath.Join(resourceFolder, meta.FileName+".md")
	if err := WriteTemplateToFile(docFileName, getDocTemplate, meta); err != nil {
		return fmt.Errorf("failed to create %s documentation: %w", kind, err)
	}

	// Update or create the `doc.go` file
	docFilePath := filepath.Join(resourceFolder, "doc.go")
	if err := updateDocFile(docFilePath, kind, resourceFolder); err != nil {
		return fmt.Errorf("failed to update 'doc.go': %w", err)
	}

	// Update the code.go file
	if err := updateLibFile(libFile, meta); err != nil {
		return fmt.Errorf("failed to update 'code.go': %w", err)
	}

	// Update the README file with a list of actions or triggers
	if err := updateReadmeFile(readmeFile, kind, meta); err != nil {
		return fmt.Errorf("failed to update 'README.md': %w", err)
	}

	fmt.Printf("%s '%s' created successfully.\n", strings.Title(kind), meta.Name)
	return nil
}

// Collects inputs interactively
func collectInput(kind string, schema *sdk.IntegrationSchemaModel, floClient *client.Client) (*ActionTriggerMetadata, error) {
	prompt := promptui.Prompt{
		Label: "Enter Name",
	}

	name, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get name: %w", err)
	}

	descMessage := fmt.Sprintf("%s integration %s called %s", schema.Name, kind, name)
	generateResponse, err := floClient.Rest.GenerateDescription(context.Background(), client.RestGenerateDescriptionRequest{
		Prompt: descMessage,
		Type:   kind,
	})
	if err != nil {
		return nil, err
	}

	descPrompt := promptui.Prompt{
		Label:     "Enter Description",
		Default:   generateResponse.Data,
		AllowEdit: true,
	}
	description, err := descPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to get description: %w", err)
	}

	// Interactive selection for type
	var typeOptions []string
	if kind == "action" {
		typeOptions = []string{
			"Normal",
		}
	} else {
		typeOptions = []string{
			"Polling",
			"Event",
			"Webhook",
			"Scheduled",
		}
	}

	typePrompt := promptui.Select{
		Label: fmt.Sprintf("Select %s Type", strings.Title(kind)),
		Items: typeOptions,
	}

	_, selectedType, err := typePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to select type: %w", err)
	}

	meta := &ActionTriggerMetadata{
		Name:        name,
		Description: strings.Trim(description, `"'`),
		Type:        selectedType,
		TypeName:    getSDKTypeName(kind, selectedType),
		FileName:    formatFileName(name),
		Constructor: getConstructorName(kind, name),
		Kind:        kind,
	}

	return meta, nil
}

func getSDKTypeName(kind, typ string) string {
	if kind == "action" {
		return fmt.Sprintf("sdkcore.ActionType%s", typ)
	}
	return fmt.Sprintf("sdkcore.TriggerType%s", typ)
}

func formatFileName(name string) string {
	// Format the name to make it file-safe (lowercase and underscores)
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}

func getConstructorName(kind, name string) string {
	return fmt.Sprintf("%s.New%s%s", kind+"s", lo.PascalCase(name), strings.Title(kind))
}

func updateLibFile(filePath string, meta *ActionTriggerMetadata) error {
	// Read the existing code.go file
	codeData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Define the insertion point based on whether it's an action or trigger
	var marker string
	if meta.Kind == "action" {
		marker = "return []sdk.Action{"
	} else {
		marker = "return []sdk.Trigger{"
	}
	insertionPoint := strings.Index(string(codeData), marker)
	if insertionPoint == -1 {
		return fmt.Errorf("failed to find the '%s' section in 'code.go'", marker)
	}

	// Prepare the line to insert
	newLine := fmt.Sprintf("\t\t%s(),", meta.Constructor)

	// Split the file content and insert the new line
	var buffer bytes.Buffer
	buffer.Write(codeData[:insertionPoint+len(marker)])
	buffer.WriteString("\n" + newLine + "\n")
	buffer.Write(codeData[insertionPoint+len(marker):])

	// Write back the updated content to code.go
	return os.WriteFile(filePath, buffer.Bytes(), 0644)
}

func updateDocFile(docFilePath, kind, resourceFolder string) error {
	// Read the existing doc.go content if it exists
	content := ""
	if _, err := os.Stat(docFilePath); err == nil {
		_, err := os.ReadFile(docFilePath)
		if err != nil {
			return fmt.Errorf("failed to read existing doc.go: %w", err)
		}
		// content = string(data)
	}

	// Find all Markdown files in the resource folder
	matches, err := filepath.Glob(filepath.Join(resourceFolder, "*.md"))
	if err != nil {
		return fmt.Errorf("failed to find markdown files: %w", err)
	}

	// Prepare the new embed declarations
	newEmbeds := ""
	for _, match := range matches {
		fileName := filepath.Base(match)
		varName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) + "Docs"
		varName = strings.ReplaceAll(varName, " ", "")
		varName = strings.ToLower(varName[:1]) + varName[1:] // Follow Go variable naming convention

		// Add the embed only if not already in the content
		if !strings.Contains(content, varName) {
			newEmbeds += fmt.Sprintf("//go:embed %s\nvar %s string\n\n", fileName, lo.CamelCase(varName))
		}
	}

	// If no new embeds are found, do nothing
	if newEmbeds == "" {
		return nil
	}

	// Ensure the package declaration exists
	if !strings.Contains(content, "package "+kind+"s") {
		content = fmt.Sprintf("package %ss\n\nimport _ \"embed\"\n\n", kind) + content
	}

	// Append new embed declarations to the content
	content += newEmbeds

	// Write back the updated content to doc.go
	if err := os.WriteFile(docFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write doc.go: %w", err)
	}

	return nil
}

func getResourceTemplate(kind string) string {
	switch kind {
	case "action":
		return actionTemplate
	case "trigger":
		return triggerTemplate
	default:
		return ""
	}
}

// Templates
const actionTemplate = `package actions

import (
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type {{ .FileName | toCamelCase }}ActionProps struct {
	Name string ` + "`json:\"name\"`" + `
}

type {{ .FileName | toPascal }}Action struct{}

func (a *{{ .FileName | toPascal }}Action) Name() string {
	return "{{ .Name }}"
}

func (a *{{ .FileName | toPascal }}Action) Description() string {
	return "{{ .Description }}"
}

func (a *{{ .FileName | toPascal }}Action) GetType() sdkcore.ActionType {
	return {{ .TypeName }}
}

func (a *{{ .FileName | toPascal }}Action) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &{{ .FileName | toCamelCase }}Docs,
	}
}

func (a *{{ .FileName | toPascal }}Action) Icon() *string {
	return nil
}

func (a *{{ .FileName | toPascal }}Action) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

func (a *{{ .FileName | toPascal }}Action) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[{{ .FileName | toCamelCase }}ActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Hello %s!", input.Name),
	}
	
	
	return out, nil
}

func (a *{{ .FileName | toPascal }}Action) Auth() *sdk.Auth {
	return nil
}

func (a *{{ .FileName | toPascal }}Action) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func (a *{{ .FileName | toPascal }}Action) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func New{{ .FileName | toPascal }}Action() sdk.Action {
	return &{{ .FileName | toPascal }}Action{}
}
`

const triggerTemplate = `package triggers

import (
	"context"
	"fmt"
	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type {{ .FileName | toCamelCase }}TriggerProps struct {
	Name string ` + "`json:\"name\"`" + `
}

type {{ .FileName | toPascal }}Trigger struct{}

func (t *{{ .FileName | toPascal }}Trigger) Name() string {
	return "{{ .Name }}"
}

func (t *{{ .FileName | toPascal }}Trigger) Description() string {
	return "{{ .Description }}"
}

func (t *{{ .FileName | toPascal }}Trigger) GetType() sdkcore.TriggerType {
	return {{ .TypeName }}
}

func (t *{{ .FileName | toPascal }}Trigger) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &{{ .FileName | toCamelCase }}Docs,
	}
}

func (t *{{ .FileName | toPascal }}Trigger) Icon() *string {
	return nil
}

func (t *{{ .FileName | toPascal }}Trigger) Properties() map[string]*sdkcore.AutoFormSchema {
	return map[string]*sdkcore.AutoFormSchema{
		"name": autoform.NewShortTextField().
			SetLabel("Name").
			SetRequired(true).
			SetPlaceholder("Your name").
			Build(),
	}
}

// Start initializes the {{ .FileName | toCamelCase }}Trigger, required for event and webhook triggers in a lifecycle context.
func (t *{{ .FileName | toPascal }}Trigger) Start(ctx sdk.LifecycleContext) error {
	// Required for event and webhook triggers
	return nil
}

// Stop shuts down the {{ .FileName | toCamelCase }}Trigger, cleaning up resources and performing necessary teardown operations.
func (t *{{ .FileName | toPascal }}Trigger) Stop(ctx sdk.LifecycleContext) error {
	return nil
}

// Execute performs the main action logic of {{ .FileName | toCamelCase }}Trigger by processing the input context and returning a JSON response.
// It converts the base context input into a strongly-typed structure, executes the desired logic, and generates output.
// Returns a JSON output map with the resulting data or an error if operation fails. required for Pooling triggers
func (t *{{ .FileName | toPascal }}Trigger) Execute(ctx sdk.ExecuteContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[{{ .FileName | toCamelCase }}TriggerProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// implement action logic
	out := map[string]any{
		"message": fmt.Sprintf("Triggered by %s!", input.Name),
	}

	return out, nil
}

func (t *{{ .FileName | toPascal }}Trigger) Criteria(ctx context.Context) sdkcore.TriggerCriteria {
	return sdkcore.TriggerCriteria{}
}

func (t *{{ .FileName | toPascal }}Trigger) Auth() *sdk.Auth {
	return nil
}

func (t *{{ .FileName | toPascal }}Trigger) SampleData() sdkcore.JSON {
	return map[string]any{
		"message": "Hello World!",
	}
}

func New{{ .FileName | toPascal }}Trigger() sdk.Trigger {
	return &{{ .FileName | toPascal }}Trigger{}
}
`

const getDocTemplate = `
# {{ .Name }}

## Description

{{ .Description }}

## Details

- **Type**: {{ .TypeName }}
`

func updateReadmeFile(readmePath, kind string, meta *ActionTriggerMetadata) error {
	var readmeContent string

	// Check if README.md exists
	if _, err := os.Stat(readmePath); err == nil {
		// Read the existing content of the README file
		content, err := os.ReadFile(readmePath)
		if err != nil {
			return fmt.Errorf("failed to read README.md file: %w", err)
		}
		readmeContent = string(content)
	} else {
		// Initialize a basic README if it doesn't exist
		readmeContent = fmt.Sprintf("# %s Integration\n\n## Description\n%s\n\n",
			meta.Name, "This integration provides various actions and triggers.")
	}

	// Define the section header for the type (actions or triggers)
	sectionHeader := fmt.Sprintf("## %s", strings.Title(kind+"s"))

	// Locate the section in the README file
	sectionStart := strings.Index(readmeContent, sectionHeader)
	var sectionContent string

	if sectionStart == -1 {
		// If section doesn't exist, initialize it with a table structure
		sectionContent = fmt.Sprintf("%s\n\n| Name | Description | Link |\n|------|-------------|------|\n", sectionHeader)
		readmeContent += sectionContent
		sectionStart = len(readmeContent) - len(sectionContent)
	} else {
		// Extract the existing section from the content
		sectionEnd := strings.Index(readmeContent[sectionStart:], "\n## ")
		if sectionEnd == -1 {
			// If this section is the last one, take everything after the header
			sectionContent = readmeContent[sectionStart:]
			readmeContent = readmeContent[:sectionStart]
		} else {
			// Extract the relevant section while keeping other content intact
			sectionEnd += sectionStart
			sectionContent = readmeContent[sectionStart:sectionEnd]
			readmeContent = readmeContent[:sectionStart] + readmeContent[sectionEnd:]
		}
	}

	// Create the new row to insert
	link := fmt.Sprintf("[docs](%s/%s.md)", kind+"s", meta.FileName)
	newRow := fmt.Sprintf("| %s | %s | %s |", meta.Name, meta.Description, link)

	// Ensure the row isn't a duplicate
	if !strings.Contains(sectionContent, link) {
		if !strings.HasSuffix(sectionContent, "\n") {
			// Ensure there is no extra newline before appending the new row
			sectionContent += "\n"
		}
		sectionContent += newRow
	}

	// Rebuild the README content
	readmeContent += sectionContent

	// Write the updated README content back to the file
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated README.md file: %w", err)
	}

	return nil
}
