package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/x-color/gaes/cmd/decrypt"
	"github.com/x-color/gaes/cmd/encrypt"
	"github.com/x-color/gaes/cmd/peek"
)

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gaes",
		Short: "gaes is tool for AES-256-CBC encryption and decryption",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(encrypt.NewEncryptCmd())
	cmd.AddCommand(decrypt.NewDecryptCmd())
	cmd.AddCommand(peek.NewPeekCmd())

	return cmd
}

func Execute() {
	cmd := rootCmd()
	cmd.SetOutput(os.Stdout)
	if err := cmd.Execute(); err != nil {
		cmd.SetOutput(os.Stderr)
		cmd.Println(err)
		os.Exit(1)
	}
}
