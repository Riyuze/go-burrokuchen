package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"time"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// Generates and returns a new block
func NewBlock(cfg *model.Config, transactions []*Transaction, prevBlockHash []byte) (*Block, error) {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}

	pow, err := NewProofOfWork(cfg, block)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	nonce, hash, err := pow.Run()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	block.Hash = hash[:]
	block.Nonce = *nonce

	return block, nil
}

// Generates and returns a genesis block
func NewGenesisBlock(cfg *model.Config, coinbase *Transaction) (*Block, error) {
	transactions := []*Transaction{coinbase}

	block, err := NewBlock(cfg, transactions, []byte{})
	if err != nil {
		return nil, utils.CatchErr(err)
	}
	return block, nil
}

func (b *Block) SerializeBlock() ([]byte, error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return result.Bytes(), nil
}

// Returns a hash of the transactions in the block
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

// Deserializes a block
func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &block, nil
}
