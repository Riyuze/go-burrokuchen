package main

import (
	"fmt"
	"go-burrokuchen/core"
	"strconv"
)

func main() {
	bc := core.NewBlockChain()

	bc.AddBlock("Send 1 BTC to Kevo")

	bc.AddBlock("Send 1 more BTC to Kevo")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)

		pow := core.NewProofOfWork(block)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
