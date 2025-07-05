package model

type Config struct {
	DatabaseConfig    DatabaseConfig
	ProofOfWorkConfig ProofOfWorkConfig
	TransactionConfig TransactionConfig
	WalletConfig      WalletConfig
	ServerConfig      ServerConfig
}

type DatabaseConfig struct {
	DbName        string
	BlocksBucket  string
	UTXOSetBucket string
}

type ProofOfWorkConfig struct {
	TargetBits int
}

type TransactionConfig struct {
	Subsidy             int
	GenesisCoinbaseData string
}

type WalletConfig struct {
	WalletFile     string
	CheckSumLength int
}

type ServerConfig struct {
	CentralNodeAddress string
	Protocol           string
	NodeVersion        int
	CommandLength      int
}
