package core

import (
	"bytes"
	"encoding/gob"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

func NewBlock(cfg *model.Config, data string, prevBlockHash []byte) (*Block, error) {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
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

func NewGenesisBlock(cfg *model.Config) (*Block, error) {
	block, err := NewBlock(cfg, "Genesis Block", []byte{})
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

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &block, nil
}
