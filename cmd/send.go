package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

func NewSendCmd(cfg *model.Config) *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Sends currency from one address to another.",
		Long:  "This command will send currency from the address that is specified to another.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := send(cfg)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	sendCmd.Flags().StringVarP(&from, "from", "f", "", "Address of the wallet sending the currency. (required)")
	sendCmd.MarkFlagRequired("from")
	sendCmd.Flags().StringVarP(&to, "to", "t", "", "Address of the wallet receiving the currency. (required)")
	sendCmd.MarkFlagRequired("to")
	sendCmd.Flags().IntVarP(&amount, "amount", "a", 0, "The amount being transferred. (required)")
	sendCmd.MarkFlagRequired("amount")

	return sendCmd
}

func send(cfg *model.Config) error {
	blockchain, err := core.InitalizeBlockchain(cfg)
	if err != nil {
		return utils.CatchErr(err)
	}
	defer blockchain.Db.Close()

	transaction, err := core.NewUTXOTransaction(*blockchain, from, to, amount)
	if err != nil {
		return utils.CatchErr(err)
	}

	err = blockchain.MineBlock([]*core.Transaction{transaction})
	if err != nil {
		return utils.CatchErr(err)
	}

	fmt.Println("Success!")

	return nil
}
