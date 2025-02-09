package commander

import "github.com/spf13/cobra"

type Command interface {
	Run(cmd *cobra.Command, args []string)
}
