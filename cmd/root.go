package cmd

import (
	"fmt"
	"github.com/wakflo/go-sdk/client"
	"log"

	"github.com/spf13/cobra"
)

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wakflo",
		Short: "Wakflo is a CLI tool for managing integrations, actions, triggers, and flows.",
		Long:  `Wakflo CLI provides functionalities to create integrations, add actions, triggers, flows, manage authentication, and more.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	floClient, err := client.New(client.Local)
	if err != nil {
		log.Fatal(err)
	}

	cmd.AddCommand(newVersionCmd(version)) // version subcommand
	//cmd.AddCommand(newExampleCmd())        // example subcommand
	cmd.AddCommand(newAuthCmd(floClient))   // auth subcommand
	cmd.AddCommand(newCreateCmd(floClient)) // create subcommand
	cmd.AddCommand(newAddCmd(floClient))    // add subcommand

	return cmd
}

// Execute invokes the command.
func Execute(version string) error {
	if err := newRootCmd(version).Execute(); err != nil {
		return fmt.Errorf("error executing root command: %w", err)
	}

	return nil
}
