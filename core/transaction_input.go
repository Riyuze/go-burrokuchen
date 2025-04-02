package core

import (
	"bytes"
	"go-burrokuchen/utils"
)

// TXInput represents a transaction input
type TXInput struct {
	TransactionID []byte
	OutputIndex   int
	Signature     []byte
	PubKey        []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) (*bool, error) {
	lockingHash, err := HashPubKey(in.PubKey)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	isKey := bytes.Equal(lockingHash, pubKeyHash)

	return &isKey, nil
}
