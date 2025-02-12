package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"

	"github.com/spf13/cobra"
	"github.com/wakflo/go-sdk/client"
	"github.com/wakflo/go-sdk/sdk"
	"github.com/wakflo/go-sdk/validator"
	"github.com/wakflo/wakflo-cli/internal/templates"
)

var val = validator.NewDefaultValidator()

func newCreateCmd(floClient *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources in Wakflo",
		Long:  "Use this command to create resources such as integrations in Wakflo.",
	}

	cmd.AddCommand(newCreateIntegrationCmd(floClient))

	return cmd
}

func newCreateIntegrationCmd(floClient *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "integration",
		Aliases: []string{"i", "int", "integ", "integrations"},
		Short:   "Create a new integration",
		Long:    "Use this command to create a new integration in Wakflo.",
		Run: func(cmd *cobra.Command, args []string) {
			// Step 1: Ask for the name of the integration
			var name string
			err := survey.AskOne(&survey.Input{
				Message: "Enter Name of the integration (required):",
			}, &name, survey.WithValidator(survey.Required))
			if err != nil {
				fmt.Printf("Name operation canceled %v+ \n", err)
				return
			}

			// Step 2: Automatically generate description
			descMessage := fmt.Sprintf("%s integration", name)
			generateResponse, err := floClient.Rest.GenerateDescription(context.Background(), client.RestGenerateDescriptionRequest{
				Prompt: descMessage,
				Type:   "integration",
			})
			if err != nil {
				fmt.Printf("Failed to generate description: %v\n", err)
				return
			}

			var description string
			err = survey.AskOne(&survey.Input{
				Message: "Enter Description of the integration (edit or accept the default):",
				Default: strings.Trim(generateResponse.Data, `"'`),
			}, &description)
			if err != nil {
				fmt.Printf("Description operation canceled %v+ \n", err)
				return
			}

			// Step 3: Fetch and choose an icon for the integration
			iconResponse, err := floClient.Rest.SearchIcon(context.Background(), client.RestSearchIconRequest{
				Name: name,
			})
			if err != nil {
				fmt.Printf("Failed to fetch icons: %v\n", err)
				return
			}

			var icon string
			if len(iconResponse.Icons) == 0 {
				err = survey.AskOne(&survey.Input{
					Message: "Enter an Icon for the integration:",
				}, &icon)
				if err != nil {
					fmt.Printf("Description operation canceled %v+ \n", err)
					return
				}
			} else {
				err = survey.AskOne(&survey.Select{
					Message: "Select an Icon for the integration:",
					Options: iconResponse.Icons,
				}, &icon)
				if err != nil {
					fmt.Printf("Icon operation canceled %v+ \n", err)
					return
				}
			}

			// Step 4: List categories and allow user to pick multiple
			catResponse, err := floClient.Rest.ListCategories(context.Background(), client.RestListCategoriesRequest{})
			if err != nil {
				fmt.Printf("Failed to fetch categories: %v\n", err)
				return
			}

			var categories []string
			err = survey.AskOne(&survey.MultiSelect{
				Message: "Select Categories for the integration:",
				Options: catResponse.Keys,
				Default: []string{"app"},
			}, &categories)
			if err != nil {
				fmt.Printf("Categories operation canceled %v+ \n", err)
				return
			}

			// Step 5: Ask for authors
			var authorsInput string
			err = survey.AskOne(&survey.Input{
				Message: "Enter Authors of the integration (comma-separated):",
				Default: "Wakflo <integrations@wakflo.com>",
			}, &authorsInput)
			if err != nil {
				fmt.Printf("Authors operation canceled %v+ \n", err)
				return
			}
			authors := strings.Split(authorsInput, ",")

			// Step 6: Automatically generate documentation
			docResponse, err := floClient.Rest.GenerateDocumentation(context.Background(), client.RestGenerateDocumentationRequest{
				Prompt: descMessage,
				Type:   "integration",
			})
			if err != nil {
				fmt.Printf("Failed to generate documentation: %v\n", err)
				return
			}

			// Step 7: Create integration metadata
			meta := &templates.CreateIntegrationProps{
				IntegrationSchemaModel: sdk.IntegrationSchemaModel{
					Name:        name,
					Description: description,
					Categories:  categories,
					Icon:        icon,
					Authors:     authors,
					Version:     "0.0.1",
				},
				Docs: docResponse.Data,
			}

			// Step 8: Validate the data
			err = val.Validate(meta)
			if err != nil {
				fmt.Printf("Invalid integration metadata: %v\n", err)
				return
			}

			// Step 9: Create the integration folder
			err = templates.CreateIntegrationFolder(meta)
			if err != nil {
				fmt.Printf("Failed to create integration: %v\n", err)
				return
			}

			fmt.Println("Integration created successfully!")
		},
	}

	return cmd
}
