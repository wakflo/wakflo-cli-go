package cmd

import (
	"github.com/FalcoSuessgott/golang-cli-template/internal/auth"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	auth := auth.New()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resources in Wakflo",
		Long:  "Use this command to create resources such as integrations in Wakflo.",
	}

	createIntegrationCmd := &cobra.Command{
		Use:   "integration",
		Short: "Create a new integration",
		Long:  "Use this command to create a new integration in Wakflo.",
		Run: func(cmd *cobra.Command, args []string) {
			auth.Login(cmd)
		},
	}

	cmd.AddCommand(createIntegrationCmd)

	return cmd
}
