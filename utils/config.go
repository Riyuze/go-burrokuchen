package utils

import (
	"go-burrokuchen/model"
	"os"
)

func LoadConfg() *model.Config {
	dbName := os.Getenv("DATABASE")
	blocksBucket := os.Getenv("BLOCKS_BUCKET")
	targetBits := os.Getenv("TARGET_BITS")
	subsidy := os.Getenv("SUBSIDY")
	genesisCoinbaseData := os.Getenv("GENESIS_COINBASE_DATA")

	cfg := &model.Config{
		DatabaseConfig: model.DatabaseConfig{
			DbName:       dbName,
			BlocksBucket: blocksBucket,
		}, ProofOfWorkConfig: model.ProofOfWorkConfig{
			TargetBits: targetBits,
		}, TransactionConfig: model.TransactionConfig{
			Subsidy:             subsidy,
			GenesisCoinbaseData: genesisCoinbaseData,
		},
	}

	return cfg
}
