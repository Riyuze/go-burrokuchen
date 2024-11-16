package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"strconv"
)

type Transaction struct {
	ID          []byte
	InputValue  []TXInput
	OutputValue []TXOutput
}

type TXInput struct {
	TransactionID []byte
	OutputIndex   int
	ScriptData    string
}

type TXOutput struct {
	Value           int
	ScriptPublicKey string
}

// Sets the ID of a transaction
func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		return utils.CatchErr(err)
	}

	hash = sha256.Sum256(encoded.Bytes())

	tx.ID = hash[:]

	return nil
}

// Generates and returns a new coinbase transaction
func NewCoinbaseTX(cfg *model.Config, to string, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Reward to: %s", to)
	}

	subsidy, err := strconv.Atoi(cfg.TransactionConfig.Subsidy)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	txIn := TXInput{TransactionID: []byte{}, OutputIndex: -1, ScriptData: data}
	txOut := TXOutput{Value: subsidy, ScriptPublicKey: to}

	tx := Transaction{ID: nil, InputValue: []TXInput{txIn}, OutputValue: []TXOutput{txOut}}
	err = tx.SetID()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &tx, nil
}

// Checks whether the transaction is coinbase or not
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.InputValue) == 1 && len(tx.InputValue[0].TransactionID) == 0 && tx.InputValue[0].OutputIndex == -1
}

// Checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptData == unlockingData
}

// Checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPublicKey == unlockingData
}

// Generates and returns a new transaction
func NewUTXOTransaction(blockchain Blockchain, from string, to string, amount int) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	balance, validOutputs, err := blockchain.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	if *balance < amount {
		err := fmt.Errorf("%s doesn't have enough funds", from)

		return nil, utils.CatchErr(err)
	}

	// Build a list of inputs
	for txID, out := range validOutputs {
		transactionID, err := hex.DecodeString(txID)
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		for _, outIndex := range out {
			input := TXInput{
				TransactionID: transactionID,
				OutputIndex:   outIndex,
				ScriptData:    from,
			}

			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	output := TXOutput{
		Value:           amount,
		ScriptPublicKey: to,
	}
	outputs = append(outputs, output)

	if *balance > amount {
		outputChange := TXOutput{
			Value:           *balance - amount,
			ScriptPublicKey: from,
		}
		outputs = append(outputs, outputChange)
	}

	tx := Transaction{
		ID:          nil,
		InputValue:  inputs,
		OutputValue: outputs,
	}
	err = tx.SetID()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &tx, nil
}
