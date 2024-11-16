package cmd

import (
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

// Flag variables
var (
	address string
	from    string
	to      string
	amount  int
)

var rootCmd = &cobra.Command{
	Use:   "go-burokkuchen",
	Short: "go-burokkuchen.",
	Long:  "go-burokkuchen is a CLI tool used for interacting with its built in blockchain written in go.",
}

func Execute() error {
	config := utils.LoadConfg()

	rootCmd.AddCommand(
		NewCreateBlockchainCmd(config),
		NewGetBalanceCmd(config),
		NewSendCmd(config),
	)

	err := rootCmd.Execute()
	if err != nil {
		return utils.CatchErr(err)
	}

	return nil
}
