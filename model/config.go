package model

type Config struct {
	DatabaseConfig    DatabaseConfig
	ProofOfWorkConfig ProofOfWorkConfig
}

type DatabaseConfig struct {
	DbName       string
	BlocksBucket string
}

type ProofOfWorkConfig struct {
	TargetBits string
}
