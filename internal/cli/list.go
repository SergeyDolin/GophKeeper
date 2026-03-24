package cli

import (
	"fmt"
	"gophkeeper/internal/clientapi"
	"gophkeeper/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "lit",
	Short: "List secrets",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()

		var password string

		fmt.Print("Master password: ")
		fmt.Scanln(&password)

		list, err := clientapi.ListSecrets(token, password)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		for _, s := range list {
			fmt.Println("-----")
			fmt.Println("Type:", s["Type"])
			fmt.Println("Data:", s["Data"])
			fmt.Println("Meta", s["Meta"])
		}
	},
}
