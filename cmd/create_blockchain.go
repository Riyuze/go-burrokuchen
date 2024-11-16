package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

func NewCreateBlockchainCmd(cfg *model.Config) *cobra.Command {
	createBlockchainCmd := &cobra.Command{
		Use:   "create-blockchain",
		Short: "Create a new blockchain",
		Long:  "This command will create a new blockchain and reward the address for mining the genesis block",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := createBlockchain(cfg)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	createBlockchainCmd.Flags().StringVarP(&address, "address", "a", "", "Address of the wallet in the blockchain. (required)")
	createBlockchainCmd.MarkFlagRequired("address")

	return createBlockchainCmd
}

func createBlockchain(cfg *model.Config) error {
	if address == "" {
		err := fmt.Errorf("please input a valid address")
		return utils.CatchErr(err)
	}

	blockchain, err := core.NewBlockchain(cfg, address)
	if err != nil {
		return utils.CatchErr(err)
	}
	defer blockchain.Db.Close()

	fmt.Println("Generated new blockhain. Sent rewards to: ", address)

	return nil
}
