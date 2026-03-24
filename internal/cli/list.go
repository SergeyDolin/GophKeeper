package cli

import (
	"encoding/base64"
	"fmt"
	"gophkeeper/internal/clientapi"
	"gophkeeper/internal/clientcrypto"
	"gophkeeper/internal/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
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

		key := clientcrypto.DeriveKey([]byte(password), []byte("static_salt"))

		for _, s := range list {
			encoded := s["Data"].(string)
			encrypted, _ := base64.StdEncoding.DecodeString(encoded)

			decrypted, err := clientcrypto.Decrypt(key, encrypted)
			if err != nil {
				fmt.Println("decrypt error")
				continue
			}
			fmt.Println("-----")
			fmt.Println("Type:", s["Type"])
			fmt.Println("Data:", string(decrypted))
			fmt.Println("Meta", s["Meta"])
		}
	},
}
