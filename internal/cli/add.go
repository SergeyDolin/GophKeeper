package cli

import (
	"fmt"
	"gophkeeper/internal/clientapi"
	"gophkeeper/internal/clientcrypto"
	"gophkeeper/internal/config"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add secret",
	Run: func(cmd *cobra.Command, args []string) {
		token, _ := config.GetToken()

		var typ, data, meta, password string

		fmt.Print("Type (text/login/card/binary):  ")
		fmt.Scanln(&typ)

		switch typ {

		case "login":
			var login, password string

			fmt.Print("Login: ")
			fmt.Scanln(&login)

			fmt.Print("Password: ")
			fmt.Scanln(&password)

			data = login + ":" + password

		case "card":
			var number, holder string

			fmt.Print("Card number: ")
			fmt.Scanln(&number)

			fmt.Print("Card holder: ")
			fmt.Scanln(&holder)

			data = number + "|" + holder

		case "text":
			fmt.Print("Text: ")
			fmt.Scanln(&data)

		default:
			fmt.Println("unknown type")
			return
		}

		fmt.Print("Meta: ")
		fmt.Scanln(&meta)

		fmt.Print("Master password: ")
		fmt.Scanln(&password)

		key := clientcrypto.DeriveKey([]byte(password), []byte("static_salt"))

		encrypted, err := clientcrypto.Encrypt(key, []byte(data))
		if err != nil {
			fmt.Println("crypto error:", err)
			return
		}

		err = clientapi.AddSecret(token, typ, meta, encrypted)
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("Encrypted & Saved")
	},
}
