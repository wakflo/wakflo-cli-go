package auth

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Auth struct {
}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) Login(cmd *cobra.Command) {
	// Add login logic
	fmt.Println("Logged in successfully!")
}

func (a *Auth) Logout(cmd *cobra.Command) {
	// Add login logic
	fmt.Println("Logged out successfully!")
}

func (a *Auth) IsLoggedIn() bool {
	return false
}

func (a *Auth) GetToken() string {
	return ""
}

func (a *Auth) SetToken(token string) {}

func (a *Auth) WhoAMI(cmd *cobra.Command) string {
	return ""
}
