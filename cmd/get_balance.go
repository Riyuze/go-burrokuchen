package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

func NewGetBalanceCmd(cfg *model.Config) *cobra.Command {
	getBalanceCmd := &cobra.Command{
		Use:   "get-balance",
		Short: "Gets the balance of an address",
		Long:  "This command will get the balance of the address that is specified",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := getBalance(cfg)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	getBalanceCmd.Flags().StringVarP(&address, "address", "a", "", "Address of the wallet in the blockchain. (required)")
	getBalanceCmd.MarkFlagRequired("address")

	return getBalanceCmd
}

func getBalance(cfg *model.Config) error {
	isValidate, err := core.ValidateAddress(cfg, address)
	if err != nil {
		return utils.CatchErr(err)
	}

	if !(*isValidate) {
		fmt.Printf("Address is not valid!")

		return nil
	}

	blockchain, err := core.InitalizeBlockchain(cfg)
	if err != nil {
		return utils.CatchErr(err)
	}
	defer blockchain.Db.Close()

	UTXSOSet := core.NewUTXOSet(cfg, blockchain)

	balance := 0

	pubKeyHash := utils.Base58Decode([]byte(address))
	checkSumLength := cfg.WalletConfig.CheckSumLength

	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumLength]

	UTXOs, err := UTXSOSet.FindUTXOByPubKeyHash(pubKeyHash)
	if err != nil {
		return utils.CatchErr(err)
	}

	for _, out := range UTXOs.Outputs {
		balance += out.Value
	}

	fmt.Printf("Balance of address '%s': %d", address, balance)

	return nil
}
