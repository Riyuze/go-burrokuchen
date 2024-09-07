package cmd

import (
	"go-burrokuchen/core"
	"go-burrokuchen/utils"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-burokkuchen",
	Short: "go-burokkuchen.",
	Long:  "go-burokkuchen is a CLI tool used for interacting with its built in blockchain written in go.",
}

func Execute() error {
	config := utils.LoadConfg()

	blockchain, err := core.NewBlockchain(config)
	if err != nil {
		log.Fatal(err)
	}

	defer blockchain.Db.Close()

	rootCmd.AddCommand(
		NewAddBlockCmd(blockchain),
		NewPrintChainCmd(blockchain),
	)

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

	return nil
}
