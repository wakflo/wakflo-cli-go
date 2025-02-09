package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add resources to Wakflo",
		Long:  "Use this command to add actions, triggers, or flows in Wakflo.",
	}

	addActionCmd := &cobra.Command{
		Use:   "action",
		Short: "Add a new action",
		Long:  "Use this command to add a new action to Wakflo.",
		Run: func(cmd *cobra.Command, args []string) {
			//auth.Login(cmd)
		},
	}

	addTriggerCmd := &cobra.Command{
		Use:   "trigger",
		Short: "Add a new trigger",
		Long:  "Use this command to add a new trigger to Wakflo.",
		Run: func(cmd *cobra.Command, args []string) {
			// Add logic to add a trigger
			fmt.Println("Trigger added successfully!")
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

	cmd.AddCommand(addTriggerCmd)
	cmd.AddCommand(addActionCmd)
	cmd.AddCommand(addFlowCmd)

	return cmd
}
