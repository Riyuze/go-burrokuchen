package utils

import (
	"go-burrokuchen/model"
	"os"
)

func LoadConfg() *model.Config {
	dbName := os.Getenv("DATABASE")
	blocksBucket := os.Getenv("BLOCKS_BUCKET")
	targetBits := os.Getenv("TARGET_BITS")

	cfg := &model.Config{
		DatabaseConfig: model.DatabaseConfig{
			DbName:       dbName,
			BlocksBucket: blocksBucket,
		}, ProofOfWorkConfig: model.ProofOfWorkConfig{
			TargetBits: targetBits,
		},
	}

	return cfg
}
