package core

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"slices"

	bolt "go.etcd.io/bbolt"
)

// Blockchain represents a blockchain
type Blockchain struct {
	cfg *model.Config
	Tip []byte
	Db  *bolt.DB
}

// NewBlockchain genearates and returns a new blockchain
func NewBlockchain(cfg *model.Config, address string) (*Blockchain, error) {
	databaseName := cfg.DatabaseConfig.DbName
	blocksBucket := []byte(cfg.DatabaseConfig.BlocksBucket)
	genesisData := cfg.TransactionConfig.GenesisCoinbaseData

	if utils.DbExists(databaseName) {
		return nil, fmt.Errorf("blockchain already exists")
	}

	var tip []byte

	db, err := bolt.Open(databaseName, 0600, nil)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		fmt.Println("No existing blockchain found. Generating a new one...")
		coinbaseTX, err := NewCoinbaseTX(cfg, address, genesisData)
		if err != nil {
			return utils.CatchErr(err)
		}

		genesis, err := NewGenesisBlock(cfg, coinbaseTX)
		if err != nil {
			return utils.CatchErr(err)
		}

		bucket, err := tx.CreateBucketIfNotExists(blocksBucket)
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

	blockChain := Blockchain{cfg: cfg, Tip: tip, Db: db}

	return &blockChain, nil
}

// InitalizeBlockchain initializes and returns a blockchain object
func InitalizeBlockchain(cfg *model.Config) (*Blockchain, error) {
	databaseName := cfg.DatabaseConfig.DbName
	blocksBucket := []byte(cfg.DatabaseConfig.BlocksBucket)

	if !utils.DbExists(databaseName) {
		return nil, fmt.Errorf("no existing blockchain found, generate one first")
	}

	var tip []byte
	db, err := bolt.Open(databaseName, 0600, nil)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)

		tip = bucket.Get([]byte("l"))

		return nil
	})

	if err != nil {
		return nil, utils.CatchErr(err)
	}

	blockchain := Blockchain{cfg: cfg, Tip: tip, Db: db}

	return &blockchain, nil
}

// MineBlock mines a new block with the provided transactions
func (bc *Blockchain) MineBlock(transactions []*Transaction) (*Block, error) {
	blocksBucket := []byte(bc.cfg.DatabaseConfig.BlocksBucket)

	var lastHash []byte

	for _, tx := range transactions {
		verified, err := bc.VerifyTransaction(tx)
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		if !*verified {
			return nil, fmt.Errorf("invalid transaction")
		}
	}

	err := bc.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)
		lastHash = bucket.Get([]byte("l"))

		return nil
	})
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	newBlock, err := NewBlock(bc.cfg, transactions, lastHash)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	err = bc.Db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(blocksBucket)
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
		return nil, utils.CatchErr(err)
	}

	return newBlock, nil
}

// InitializeIterator initializes the blockchain iterator object
func (bc *Blockchain) InitializeIterator() *BlockchainIterator {
	bci := &BlockchainIterator{cfg: bc.cfg, currentHash: bc.Tip, db: bc.Db}

	return bci
}

// FindUTXO finds all unspent transaction outputs and returns transactions with spent outputs removed
func (bc *Blockchain) FindUTXO() (map[string]TXOutputs, error) {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)

	bci := bc.InitializeIterator()

	for {
		block, err := bci.Prev()
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		for _, transaction := range block.Transactions {
			transactionID := hex.EncodeToString(transaction.ID)

			err := func() error {
				for outIndex, out := range transaction.OutputValue {
					if spentTXOs[transactionID] != nil {
						if slices.Contains(spentTXOs[transactionID], outIndex) {
							return nil
						}
					}

					outs := UTXO[transactionID]
					outs.Outputs = append(outs.Outputs, out)
					UTXO[transactionID] = outs
				}

				if !transaction.IsCoinbase() {
					for _, in := range transaction.InputValue {
						inTransactionID := hex.EncodeToString(in.TransactionID)
						spentTXOs[inTransactionID] = append(spentTXOs[inTransactionID], in.OutputIndex)
					}
				}

				return nil
			}()

			if err != nil {
				return nil, utils.CatchErr(err)
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO, nil
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (*Transaction, error) {
	bci := bc.InitializeIterator()

	for {
		block, err := bci.Prev()
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return &Transaction{}, fmt.Errorf("transaction not found")
}

// SignTransaction signs inputs of a Transaction
func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) error {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.InputValue {
		prevTX, err := bc.FindTransaction(vin.TransactionID)
		if err != nil {
			return utils.CatchErr(err)
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = *prevTX
	}

	tx.Sign(privKey, prevTXs)

	return nil
}

// VerifyTransaction verifies transaction input signatures
func (bc *Blockchain) VerifyTransaction(tx *Transaction) (*bool, error) {
	if tx.IsCoinbase() {
		verified := true

		return &verified, nil
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.InputValue {
		prevTX, err := bc.FindTransaction(vin.TransactionID)
		if err != nil {
			return nil, utils.CatchErr(err)
		}

		prevTXs[hex.EncodeToString(prevTX.ID)] = *prevTX
	}

	verified, err := tx.Verify(prevTXs)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return verified, nil
}
