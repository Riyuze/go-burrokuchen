package core

import (
	"encoding/hex"
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

// Genearates and returns a new blockchain
func NewBlockchain(cfg *model.Config, address string) (*Blockchain, error) {
	if utils.DbExists(cfg.DatabaseConfig.DbName) {
		return nil, fmt.Errorf("blockchain already exists")
	}

	var tip []byte

	db, err := bolt.Open(cfg.DatabaseConfig.DbName, 0600, nil)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		fmt.Println("No existing blockchain found. Generating a new one...")
		coinbaseTX, err := NewCoinbaseTX(cfg, address, cfg.TransactionConfig.GenesisCoinbaseData)
		if err != nil {
			return utils.CatchErr(err)
		}

		genesis, err := NewGenesisBlock(cfg, coinbaseTX)
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

		return nil
	})

	if err != nil {
		return nil, utils.CatchErr(err)
	}

	blockChain := Blockchain{Cfg: cfg, Tip: tip, Db: db}

	return &blockChain, nil
}

// Initializes and returns a blockchain object
func InitalizeBlockchain(cfg *model.Config) (*Blockchain, error) {
	if !utils.DbExists(cfg.DatabaseConfig.DbName) {
		return nil, fmt.Errorf("no existing blockchain found, generate one first")
	}

	var tip []byte
	db, err := bolt.Open(cfg.DatabaseConfig.DbName, 0600, nil)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(cfg.DatabaseConfig.BlocksBucket))

		tip = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		return nil, utils.CatchErr(err)
	}

	blockchain := Blockchain{Cfg: cfg, Tip: tip, Db: db}

	return &blockchain, nil
}

// Mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) error {
	var lastHash []byte

	err := bc.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bc.Cfg.DatabaseConfig.BlocksBucket))
		lastHash = bucket.Get([]byte("l"))

		return nil
	})
	if err != nil {
		return utils.CatchErr(err)
	}

	newBlock, err := NewBlock(bc.Cfg, transactions, lastHash)
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

// Initializes the blockchain iterator object
func (bc *Blockchain) InitializeIterator() *BlockchainIterator {
	bci := &BlockchainIterator{cfg: bc.Cfg, currentHash: bc.Tip, db: bc.Db}

	return bci
}

// Returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(address string) ([]*Transaction, error) {
	var unspentTXs []*Transaction
	spentTXOs := make(map[string][]int)

	bci := bc.InitializeIterator()

	for {
		block, err := bci.Prev()
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

			func() {
				for outIndex, out := range transaction.OutputValue {
					if spentTXOs[transactionID] != nil {
						for _, spentOut := range spentTXOs[transactionID] {
							if spentOut == outIndex {
								return
							}
						}
					}

					if out.CanBeUnlockedWith(address) {
						unspentTXs = append(unspentTXs, transaction)
					}
				}

				if !transaction.IsCoinbase() {
					for _, in := range transaction.InputValue {
						if in.CanUnlockOutputWith(address) {
							inTransactionID := hex.EncodeToString(in.TransactionID)
							spentTXOs[inTransactionID] = append(spentTXOs[inTransactionID], in.OutputIndex)
						}
					}
				}
			}()
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs, nil
}

// Finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUnspentTransactionOutputs(address string) ([]*TXOutput, error) {
	var unspentTXOs []*TXOutput
	unspentTXs, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	for _, transaction := range unspentTXs {
		for _, out := range transaction.OutputValue {
			if out.CanBeUnlockedWith(address) {
				unspentTXOs = append(unspentTXOs, &out)
			}
		}
	}

	return unspentTXOs, nil
}

// Finds and returns unspent outputs in reference to an amount
func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (*int, map[string][]int, error) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	unspentTXs, err := bc.FindUnspentTransactions(address)
	if err != nil {
		return nil, nil, utils.CatchErr(err)
	}

	func() {
		for _, transaction := range unspentTXs {
			transactionID := hex.EncodeToString(transaction.ID)

			for outIndex, out := range transaction.OutputValue {
				if out.CanBeUnlockedWith(address) && accumulated < amount {
					accumulated += out.Value

					unspentOutputs[transactionID] = append(unspentOutputs[transactionID], outIndex)

					if accumulated >= amount {
						return
					}
				}
			}
		}
	}()

	return &accumulated, unspentOutputs, nil
}
