package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"math/big"
)

// Transaction represents a transaction
type Transaction struct {
	ID          []byte
	InputValue  []TXInput
	OutputValue []TXOutput
}

// Serialize returns a serialized Transaction
func (tx Transaction) Serialize() ([]byte, error) {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return encoded.Bytes(), nil
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() ([]byte, error) {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	serializedTx, err := txCopy.Serialize()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	hash = sha256.Sum256(serializedTx)

	return hash[:], nil
}

// NewCoinbaseTX generates and returns a new coinbase transaction
func NewCoinbaseTX(cfg *model.Config, to string, data string) (*Transaction, error) {
	subsidy := cfg.TransactionConfig.Subsidy

	if data == "" {
		data = fmt.Sprintf("Reward sent to: %s", to)
	}

	txIn := TXInput{
		TransactionID: []byte{},
		OutputIndex:   -1,
		Signature:     nil,
		PubKey:        []byte(data),
	}
	txOut, err := NewTXOutput(cfg, subsidy, to)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	tx := Transaction{ID: nil, InputValue: []TXInput{txIn}, OutputValue: []TXOutput{*txOut}}
	hash, err := tx.Hash()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	tx.ID = hash

	return &tx, nil
}

// IsCoinbase checks whether the transaction is coinbase or not
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.InputValue) == 1 && len(tx.InputValue[0].TransactionID) == 0 && tx.InputValue[0].OutputIndex == -1
}

// NewUTXOTransaction generates and returns a new transaction
func NewUTXOTransaction(utxoSet UTXOSet, from string, to string, amount int) (*Transaction, error) {
	var inputs []TXInput
	var outputs []TXOutput

	wallets, err := NewWallets(utxoSet.cfg)
	if err != nil {
		return nil, utils.CatchErr(err)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash, err := HashPubKey(wallet.PublicKey)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	balance, validOutputs, err := utxoSet.FindSpendableOutputs(pubKeyHash, amount)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	if *balance < amount {
		err := fmt.Errorf("%s doesn't have enough funds", from)

		return nil, utils.CatchErr(err)
	}

	// Build a list of inputs
	for txID, outputs := range validOutputs {
		transactionID, err := hex.DecodeString(txID)
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		for _, outIndex := range outputs {
			input := TXInput{
				TransactionID: transactionID,
				OutputIndex:   outIndex,
				Signature:     nil,
				PubKey:        wallet.PublicKey,
			}

			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	output, err := NewTXOutput(utxoSet.cfg, amount, to)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	outputs = append(outputs, *output)

	if *balance > amount {
		outputChange, err := NewTXOutput(utxoSet.cfg, *balance-amount, from)
		if err != nil {
			return nil, utils.CatchErr(err)
		}
		outputs = append(outputs, *outputChange)
	}

	tx := Transaction{
		ID:          nil,
		InputValue:  inputs,
		OutputValue: outputs,
	}
	hash, err := tx.Hash()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	tx.ID = hash

	utxoSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx, nil
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.InputValue {
		inputs = append(inputs, TXInput{TransactionID: vin.TransactionID, OutputIndex: vin.OutputIndex, Signature: nil, PubKey: nil})
	}

	for _, vout := range tx.OutputValue {
		outputs = append(outputs, TXOutput{Value: vout.Value, PubKeyHash: vout.PubKeyHash})
	}

	txCopy := Transaction{ID: tx.ID, InputValue: inputs, OutputValue: outputs}

	return txCopy
}

// Sign signs each input of a Transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	txCopy := tx.TrimmedCopy()

	for inputIndex, vin := range txCopy.InputValue {
		prevTX := prevTXs[hex.EncodeToString(vin.TransactionID)]
		txCopy.InputValue[inputIndex].Signature = nil
		txCopy.InputValue[inputIndex].PubKey = prevTX.OutputValue[vin.OutputIndex].PubKeyHash

		hashValue, err := txCopy.Hash()
		if err != nil {
			return utils.CatchErr(err)
		}
		txCopy.ID = hashValue
		txCopy.InputValue[inputIndex].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if err != nil {
			return utils.CatchErr(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.InputValue[inputIndex].Signature = signature

	}

	return nil
}

// Verify verifies signatures of Transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) (*bool, error) {
	var verified bool

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inputIndex, vin := range tx.InputValue {
		prevTX := prevTXs[hex.EncodeToString(vin.TransactionID)]
		txCopy.InputValue[inputIndex].Signature = nil
		txCopy.InputValue[inputIndex].PubKey = prevTX.OutputValue[vin.OutputIndex].PubKeyHash

		hashValue, err := txCopy.Hash()
		if err != nil {
			return &verified, utils.CatchErr(err)
		}
		txCopy.ID = hashValue
		txCopy.InputValue[inputIndex].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return &verified, nil
		}
	}

	verified = true

	return &verified, nil
}
