package encrypt

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/x-color/gaes/crypto"
	"github.com/x-color/gaes/path"
)

func NewEncryptCmd() *cobra.Command {
	var isOverWrite bool
	var password string

	cmd := &cobra.Command{
		Use:     "encrypt in [out]",
		Aliases: []string{"e"},
		Short:   "encrypt plain file",
		Example: `  gaes encrypto plain.txt encrypted.txt
  gaes e plain.txt encrypted.txt
  gaes e -w plain.txt`,
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

			return encryptFile(key, in, out)
		},
	}

	cmd.Flags().BoolVarP(&isOverWrite, "overwrite", "w", false, "overwrite mode")
	cmd.Flags().StringVarP(&password, "password", "p", "", "password")

	return cmd
}

func encryptFile(key []byte, fromPath, toPath string) error {
	plainText, err := ioutil.ReadFile(fromPath)
	if err != nil {
		return err
	}

	chipherText, err := crypto.Encrypt(key, plainText)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(toPath, chipherText, 0644)
}

func readPwd(cmd *cobra.Command) ([]byte, error) {
	cmd.Print("enter the password: ")
	defer cmd.Println()
	pwd, err := crypto.ReadPwd()
	if err != nil {
		return nil, err
	}
	cmd.Println()

	cmd.Print("valify the password: ")
	valify, err := crypto.ReadPwd()
	if err != nil {
		return nil, err
	}

	if string(pwd) != string(valify) {
		return nil, errors.New("password valify failed")
	}

	return pwd, nil
}
