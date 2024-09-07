package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/utils"

	"github.com/spf13/cobra"
)

var data string

func NewAddBlockCmd(blockchain *core.Blockchain) *cobra.Command {
	addBlockCmd := &cobra.Command{
		Use:   "add-block",
		Short: "Add a block to the blockchain.",
		Long:  "This command will add a block into the blockchain. If a blockchain is not found, it will generate a new one.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := addBlock(blockchain)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	addBlockCmd.Flags().StringVarP(&data, "data", "d", "", "Block data to be added (required)")
	addBlockCmd.MarkFlagRequired("data")

	return addBlockCmd
}

func addBlock(blockchain *core.Blockchain) error {
	err := blockchain.AddBlock(data)
	if err != nil {
		return utils.CatchErr(err)
	}

	fmt.Println("Success!")

	return nil
}
