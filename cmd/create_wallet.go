package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

func NewCreateWalletCmd(cfg *model.Config) *cobra.Command {
	createWalletCmd := &cobra.Command{
		Use:   "create-wallet",
		Short: "Creates a new wallet",
		Long:  "This command will create a new wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := createWallet(cfg)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	return createWalletCmd
}

func createWallet(cfg *model.Config) error {
	wallets, err := core.NewWallets(cfg)
	if err != nil {
		return utils.CatchErr(err)
	}

	address, err := wallets.CreateWallet()
	if err != nil {
		return utils.CatchErr(err)
	}

	err = wallets.SaveToFile()
	if err != nil {
		return utils.CatchErr(err)
	}

	fmt.Printf("Your new address: %s", *address)

	return nil
}
