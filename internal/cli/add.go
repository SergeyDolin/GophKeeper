package cli

import (
	"fmt"
	"gophkeeper/internal/clientapi"
	"gophkeeper/internal/config"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add secret",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()

		var typ, data, meta, password string

		fmt.Print("Type: ")
		fmt.Scanln(&typ)

		fmt.Print("Data: ")
		fmt.Scanln(&data)

		fmt.Print("Meta: ")
		fmt.Scanln(&meta)

		fmt.Print("Master password: ")
		fmt.Scanln(&password)

		err := clientapi.AddSecret(token, typ, data, meta, password)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("Saved")
	},
}
