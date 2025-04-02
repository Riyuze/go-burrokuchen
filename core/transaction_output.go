package core

import (
	"bytes"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"strconv"
)

// TXOutput represents a transaction output
type TXOutput struct {
	cfg        *model.Config
	Value      int
	PubKeyHash []byte
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) error {
	checkSumLength, err := strconv.Atoi(out.cfg.WalletConfig.CheckSumLength)
	if err != nil {
		return utils.CatchErr(err)
	}

	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumLength]
	out.PubKeyHash = pubKeyHash

	return nil
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Equal(out.PubKeyHash, pubKeyHash)
}

// NewTXOutput create a new TXOutput
func NewTXOutput(cfg *model.Config, value int, address string) (*TXOutput, error) {
	txo := &TXOutput{cfg: cfg, Value: value, PubKeyHash: nil}
	err := txo.Lock([]byte(address))
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return txo, nil
}
