package core

import (
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	bolt "go.etcd.io/bbolt"
)

type Blockchain struct {
	Cfg *model.Config
	Tip []byte
	Db  *bolt.DB
}

func NewBlockchain(cfg *model.Config) (*Blockchain, error) {
	var tip []byte

	db, err := bolt.Open(cfg.DatabaseConfig.DbName, 0600, nil)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(cfg.DatabaseConfig.BlocksBucket))

		if bucket == nil {
			fmt.Println("No existing blockchain found. Generating a new one...")
			genesis, err := NewGenesisBlock(cfg)
			if err != nil {
				return utils.CatchErr(err)
			}

			bucket, err := tx.CreateBucket([]byte(cfg.DatabaseConfig.BlocksBucket))
			if err != nil {
				return utils.CatchErr(err)
			}

			serializedBlock, err := genesis.SerializeBlock()
			if err != nil {
				return utils.CatchErr(err)
			}

			err = bucket.Put(genesis.Hash, serializedBlock)
			if err != nil {
				return utils.CatchErr(err)
			}

			err = bucket.Put([]byte("l"), genesis.Hash)
			if err != nil {
				return utils.CatchErr(err)
			}

			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte("l"))
		}

		return nil
	})

	if err != nil {
		return nil, utils.CatchErr(err)
	}

	blockChain := Blockchain{Cfg: cfg, Tip: tip, Db: db}

	return &blockChain, nil
}

func (bc *Blockchain) AddBlock(data string) error {
	var lastHash []byte

	err := bc.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bc.Cfg.DatabaseConfig.BlocksBucket))
		lastHash = bucket.Get([]byte("l"))

		return nil
	})
	if err != nil {
		return utils.CatchErr(err)
	}

	newBlock, err := NewBlock(bc.Cfg, data, lastHash)
	if err != nil {
		return utils.CatchErr(err)
	}

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bc.Cfg.DatabaseConfig.BlocksBucket))
		serializedBlock, err := newBlock.SerializeBlock()
		if err != nil {
			return utils.CatchErr(err)
		}

		err = bucket.Put(newBlock.Hash, serializedBlock)
		if err != nil {
			return utils.CatchErr(err)
		}

		err = bucket.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			return utils.CatchErr(err)
		}

		bc.Tip = newBlock.Hash

		return nil
	})
	if err != nil {
		return utils.CatchErr(err)
	}

	return nil
}

func (bc *Blockchain) InitializeIterator() *BlockchainIterator {
	bci := &BlockchainIterator{cfg: bc.Cfg, currentHash: bc.Tip, db: bc.Db}

	return bci
}
