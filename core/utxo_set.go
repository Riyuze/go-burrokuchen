package core

import (
	"encoding/hex"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"

	"go.etcd.io/bbolt"
)

// UTXOSet represents UTXO set
type UTXOSet struct {
	cfg        *model.Config
	Blockchain *Blockchain
}

// NewUTXOSet generates and returns a new UTXO set object
func NewUTXOSet(cfg *model.Config, blockchain *Blockchain) *UTXOSet {
	return &UTXOSet{cfg: cfg, Blockchain: blockchain}
}

// Reindex rebuilds the UTXO set
func (u *UTXOSet) Reindex() error {
	utxoSetBucket := []byte(u.cfg.DatabaseConfig.UTXOSetBucket)
	db := u.Blockchain.Db

	err := db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket(utxoSetBucket)
		if err != nil && err != bbolt.ErrBucketNotFound {
			return utils.CatchErr(err)
		}

		_, err = tx.CreateBucket(utxoSetBucket)
		if err != nil {
			return utils.CatchErr(err)
		}

		return nil
	})

	if err != nil {
		return utils.CatchErr(err)
	}

	UTXO, err := u.Blockchain.FindUTXO()
	if err != nil {
		return utils.CatchErr(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(utxoSetBucket)

		for txID, outputs := range UTXO {
			key, err := hex.DecodeString(txID)
			if err != nil {
				return utils.CatchErr(err)
			}

			serializedOutputs, err := outputs.Serialize()
			if err != nil {
				return utils.CatchErr(err)
			}

			err = b.Put(key, serializedOutputs)
			if err != nil {
				return utils.CatchErr(err)
			}
		}

		return nil
	})

	if err != nil {
		return utils.CatchErr(err)
	}

	return nil
}

// Update updates the UTXO set with transactions from the Block (tip of the blockchain)
func (u UTXOSet) Update(block *Block) error {
	utxoSetBucket := []byte(u.cfg.DatabaseConfig.UTXOSetBucket)
	db := u.Blockchain.Db

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(utxoSetBucket)

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, vin := range tx.InputValue {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.TransactionID)
					outs, err := DeserializeOutputs(outsBytes)
					if err != nil {
						return utils.CatchErr(err)
					}

					for outIndex, out := range outs.Outputs {
						if outIndex != vin.OutputIndex {
							updatedOuts.Outputs = append(updatedOuts.Outputs, out)
						}
					}

					if len(updatedOuts.Outputs) == 0 {
						err := b.Delete(vin.TransactionID)
						if err != nil {
							return utils.CatchErr(err)
						}
					} else {
						serializedOutputs, err := updatedOuts.Serialize()
						if err != nil {
							return utils.CatchErr(err)
						}

						err = b.Put(vin.TransactionID, serializedOutputs)
						if err != nil {
							return utils.CatchErr(err)
						}
					}
				}
			}

			newOutputs := TXOutputs{}
			newOutputs.Outputs = append(newOutputs.Outputs, tx.OutputValue...)

			serializedOutputs, err := newOutputs.Serialize()
			if err != nil {
				return utils.CatchErr(err)
			}

			err = b.Put(tx.ID, serializedOutputs)
			if err != nil {
				return utils.CatchErr(err)
			}
		}

		return nil
	})

	if err != nil {
		return utils.CatchErr(err)
	}

	return nil
}

// FindSpendableOutputs finds and returns unspent outputs in reference to an amount
func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (*int, map[string][]int, error) {
	utxoSetBucket := []byte(u.cfg.DatabaseConfig.UTXOSetBucket)
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	db := u.Blockchain.Db

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(utxoSetBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs, err := DeserializeOutputs(v)
			if err != nil {
				return utils.CatchErr(err)
			}

			for outIndex, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIndex)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, nil, utils.CatchErr(err)
	}

	return &accumulated, unspentOutputs, nil
}

// FindUTXOByPubKeyHash finds UTXO for a public key hash
func (u *UTXOSet) FindUTXOByPubKeyHash(pubKeyHash []byte) (*TXOutputs, error) {
	utxoSetBucket := []byte(u.cfg.DatabaseConfig.UTXOSetBucket)
	var UTXOs TXOutputs
	db := u.Blockchain.Db

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(utxoSetBucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs, err := DeserializeOutputs(v)
			if err != nil {
				return utils.CatchErr(err)
			}

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs.Outputs = append(UTXOs.Outputs, out)
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &UTXOs, nil
}
