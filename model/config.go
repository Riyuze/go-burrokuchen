package model

type Config struct {
	DatabaseConfig    DatabaseConfig
	ProofOfWorkConfig ProofOfWorkConfig
	TransactionConfig TransactionConfig
}

type DatabaseConfig struct {
	DbName       string
	BlocksBucket string
}

type ProofOfWorkConfig struct {
	TargetBits string
}

type TransactionConfig struct {
	Subsidy             string
	GenesisCoinbaseData string
}
