package decrypt

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/x-color/gaes/crypto"
	"github.com/x-color/gaes/path"
)

func NewDecryptCmd() *cobra.Command {
	var isOverWrite bool
	var password string

	cmd := &cobra.Command{
		Use:     "decrypt in [out]",
		Aliases: []string{"d"},
		Short:   "decrypt encrypted file",
		Example: `  gaes decrypto encrypted.txt plain.txt
  gaes d encrypted.txt plain.txt
  gaes d -w encrypted.txt`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if isOverWrite {
				if len(args) == 2 {
					return errors.New("do not need second argument(out) in overwrite mode")
				}
			} else if len(args) == 1 {
				return errors.New("need second argument(out)")
			}

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

			in, err := path.Abs(args[0])
			if err != nil {
				return err
			}

			var second string
			if isOverWrite {
				second = args[0]
			} else {
				second = args[1]
			}

			out, err := path.Abs(second)
			if err != nil {
				return err
			}

			return decryptFile(key, in, out)
		},
	}

	cmd.Flags().BoolVarP(&isOverWrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password")

	return cmd
}

func decryptFile(key []byte, fromPath, toPath string) error {
	chipherText, err := ioutil.ReadFile(fromPath)
	if err != nil {
		return err
	}

	plainText, err := crypto.Decrypt(key, chipherText)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(toPath, plainText, 0644)
}

func readPwd(cmd *cobra.Command) ([]byte, error) {
	cmd.Print("enter the password: ")
	defer cmd.Println()
	return crypto.ReadPwd()
}
