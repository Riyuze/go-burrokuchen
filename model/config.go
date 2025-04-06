package model

type Config struct {
	DatabaseConfig    DatabaseConfig
	ProofOfWorkConfig ProofOfWorkConfig
	TransactionConfig TransactionConfig
	WalletConfig      WalletConfig
}

type DatabaseConfig struct {
	DbName        string
	BlocksBucket  string
	UTXOSetBucket string
}

type ProofOfWorkConfig struct {
	TargetBits string
}

type TransactionConfig struct {
	Subsidy             string
	GenesisCoinbaseData string
}

type WalletConfig struct {
	WalletFile     string
	CheckSumLength string
}
