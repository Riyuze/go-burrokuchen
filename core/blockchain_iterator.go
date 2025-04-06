package core

import (
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	bolt "go.etcd.io/bbolt"
)

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	cfg         *model.Config
	currentHash []byte
	db          *bolt.DB
}

// NewBlockchainIterator generates and returns a blockchain iterator
func NewBlockchainIterator(cfg *model.Config, tip []byte, db *bolt.DB) *BlockchainIterator {
	return &BlockchainIterator{cfg: cfg, currentHash: tip, db: db}
}

// Prev returns the previous block instance in the blockchain
func (bci *BlockchainIterator) Prev() (*Block, error) {
	var prevBlock *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bci.cfg.DatabaseConfig.BlocksBucket))
		encodedBlock := bucket.Get(bci.currentHash)
		block, err := DeserializeBlock(encodedBlock)
		if err != nil {
			return utils.CatchErr(err)
		}

		prevBlock = block

		return nil
	})
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	bci.currentHash = prevBlock.PrevBlockHash

	return prevBlock, nil
}
