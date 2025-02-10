package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wakflo/go-sdk/client"
	"github.com/wakflo/wakflo-cli/internal/auth"
)

func newAuthCmd(floClient *client.Client) *cobra.Command {
	auth := auth.New()

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication for Wakflo",
		Long:  "Use this command to log in or log out of Wakflo.",
	}

	authLoginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Wakflo",
		Long:  "Use this command to log in to Wakflo and authenticate your session.",
		Run: func(cmd *cobra.Command, args []string) {
			auth.Login(cmd)
		},
	}

	authLogoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of Wakflo",
		Long:  "Use this command to log out of Wakflo and end your session.",
		Run: func(cmd *cobra.Command, args []string) {
			auth.Logout(cmd)
		},
	}

	cmd.AddCommand(authLoginCmd)
	cmd.AddCommand(authLogoutCmd)

	return cmd
}
