package cmd

import (
	"fmt"
	"go-burrokuchen/core"
	"go-burrokuchen/utils"
	"strconv"

	"github.com/spf13/cobra"
)

func NewPrintChainCmd(blockchain *core.Blockchain) *cobra.Command {
	printChainCmd := &cobra.Command{
		Use:   "print-chain",
		Short: "Prints the blockchain",
		Long:  "This command will iterate through the blockchain, printing each of it's blocks contents.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := printChain(blockchain)
			if err != nil {
				return utils.CatchErr(err)
			}

			return nil
		},
	}

	return printChainCmd
}

func printChain(blockchain *core.Blockchain) error {
	bci := blockchain.InitializeIterator()

	for {
		block, err := bci.Prev()
		if err != nil {
			return utils.CatchErr(err)
		}

		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow, err := core.NewProofOfWork(blockchain.Cfg, block)
		if err != nil {
			return utils.CatchErr(err)
		}

		validation, err := pow.Validate()
		if err != nil {
			return utils.CatchErr(err)
		}

		fmt.Printf("PoW validated: %s\n", strconv.FormatBool(*validation))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return nil
}
