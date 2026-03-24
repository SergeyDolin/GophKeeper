package cli

import (
	"fmt"
	"gophkeeper/internal/clientapi"
	"gophkeeper/internal/config"

	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login user",
	Run: func(cmd *cobra.Command, args []string) {
		var login, password string

		fmt.Print("Login: ")
		fmt.Scanln(&login)

		fmt.Print("Password: ")
		fmt.Scanln(&password)

		token, err := clientapi.Login(login, password)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		config.SaveToken(token)

		fmt.Println("Logged in")
	},
}
