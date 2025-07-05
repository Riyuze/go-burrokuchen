package utils

import (
	"go-burrokuchen/model"

	"github.com/spf13/viper"
)

func LoadConfg() (*model.Config, error) {
	vip := viper.New()
	vip.SetConfigName("config.yaml")
	vip.SetConfigType("yaml")

	vip.AddConfigPath(".")
	err := vip.ReadInConfig()
	if err != nil {
		return nil, CatchErr(err)
	}

	dbName := vip.GetString("database.name")
	blocksBucket := vip.GetString("database.blocks_bucket")
	utxoSetBucket := vip.GetString("database.utxo_set_bucket")
	targetBits := vip.GetInt("proof_of_work.target_bits")
	subsidy := vip.GetInt("transaction.subsidy")
	genesisCoinbaseData := vip.GetString("transaction.genesis_coinbase_data")
	walletFile := vip.GetString("wallet.file")
	checkSumLength := vip.GetInt("wallet.check_sum_length")
	centralNodeAddress := vip.GetString("server.central_node")
	protocol := vip.GetString("server.protocol")
	nodeVersion := vip.GetInt("server.node_version")
	commandLength := vip.GetInt("server.command_length")

	cfg := &model.Config{
		DatabaseConfig: model.DatabaseConfig{
			DbName:        dbName,
			BlocksBucket:  blocksBucket,
			UTXOSetBucket: utxoSetBucket,
		}, ProofOfWorkConfig: model.ProofOfWorkConfig{
			TargetBits: targetBits,
		}, TransactionConfig: model.TransactionConfig{
			Subsidy:             subsidy,
			GenesisCoinbaseData: genesisCoinbaseData,
		}, WalletConfig: model.WalletConfig{
			WalletFile:     walletFile,
			CheckSumLength: checkSumLength,
		}, ServerConfig: model.ServerConfig{
			CentralNodeAddress: centralNodeAddress,
			Protocol:           protocol,
			NodeVersion:        nodeVersion,
			CommandLength:      commandLength,
		},
	}

	return cfg, nil
}
