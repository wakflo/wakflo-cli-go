package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wakflo/go-sdk/client"
	"github.com/wakflo/wakflo-cli/internal/templates"
)

func newAddCmd(floClient *client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add resources to Wakflo",
		Long:  "Use this command to add actions, triggers, or flows in Wakflo.",
	}

	for _, operationCmd := range newOperationsCmd(floClient) {
		cmd.AddCommand(operationCmd)
	}

	return cmd
}

func newOperationsCmd(floClient *client.Client) []*cobra.Command {
	// Subcommand for adding an action
	addActionCmd := &cobra.Command{
		Use:   "action",
		Short: "Add a new action to the integration",
		Long:  "Use this command to add a new action to the current integration project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return templates.HandleAddResource("action", cmd, floClient)
		},
	}

	// Subcommand for adding a trigger
	addTriggerCmd := &cobra.Command{
		Use:   "trigger",
		Short: "Add a new trigger to the integration",
		Long:  "Use this command to add a new trigger to the current integration project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return templates.HandleAddResource("trigger", cmd, floClient)
		},
	}

	addFlowCmd := &cobra.Command{
		Use:   "flow",
		Short: "Add a new flow",
		Long:  "Use this command to add a new flow to Wakflo.",
		Run: func(cmd *cobra.Command, args []string) {
			// Add logic to add a flow
			fmt.Println("Flow added successfully!")
		},
	}

	return []*cobra.Command{addActionCmd, addTriggerCmd, addFlowCmd}
}

func registerAddFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("name", "n", "", "Name of the resource (required)")
	cmd.Flags().StringP("description", "d", "", "Description of the resource")
	cmd.Flags().StringP("type", "t", "", "Type of the resource (e.g., sdkcore.ActionType or sdkcore.TriggerType)")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("type")
}
