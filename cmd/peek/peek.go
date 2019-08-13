package peek

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/x-color/gaes/crypto"
	"github.com/x-color/gaes/path"
)

func NewPeekCmd() *cobra.Command {
	var password string

	cmd := &cobra.Command{
		Use:     "peek file",
		Aliases: []string{"p"},
		Short:   "peek plain text in encrypted file",
		Example: `  gaes peek encrypted.txt
  gaes p encrypted.txt`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var key []byte
			if password == "" {
				pwd, err := readPwd(cmd)
				if err != nil {
					return err
				}
				key = crypto.PwdToKey(pwd)
			} else {
				key = crypto.PwdToKey([]byte(password))
			}

			file, err := path.Abs(args[0])
			if err != nil {
				return err
			}

			plainText, err := peekFile(key, file)
			if err != nil {
				return err
			}
			cmd.Print(string(plainText))

			return nil
		},
	}

	cmd.Flags().StringVarP(&password, "password", "p", "", "password")

	return cmd
}

func peekFile(key []byte, filepath string) ([]byte, error) {
	chipherText, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return crypto.Decrypt(key, chipherText)
}

func readPwd(cmd *cobra.Command) ([]byte, error) {
	cmd.Print("enter the password: ")
	defer cmd.Println()
	return crypto.ReadPwd()
}
