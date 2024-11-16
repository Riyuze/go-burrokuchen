package core

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"math"
	"math/big"
	"strconv"
)

var (
	maxNonce = math.MaxInt64
)

type ProofOfWork struct {
	cfg    *model.Config
	block  *Block
	target *big.Int
}

// Generates and returns a proof of work
func NewProofOfWork(cfg *model.Config, b *Block) (*ProofOfWork, error) {
	targetBits, err := strconv.Atoi(cfg.ProofOfWorkConfig.TargetBits)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{
		cfg:    cfg,
		block:  b,
		target: target,
	}

	return pow, nil
}

// Prepares data for the proof of work
func (pow *ProofOfWork) prepareData(nonce int) ([]byte, error) {
	targetBits, err := strconv.Atoi(pow.cfg.ProofOfWorkConfig.TargetBits)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	timestampBytes, err := utils.IntToHex(pow.block.Timestamp)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	targetBitsBytes, err := utils.IntToHex(int64(targetBits))
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	nonceBytes, err := utils.IntToHex(int64(nonce))
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			timestampBytes,
			targetBitsBytes,
			nonceBytes,
		}, []byte{},
	)

	return data, nil
}

// Runs the proof of work
func (pow *ProofOfWork) Run() (*int, []byte, error) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Println("Mining a new block...")
	for nonce < maxNonce {
		data, err := pow.prepareData(nonce)
		if err != nil {
			return nil, nil, utils.CatchErr(err)
		}

		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce += 1
		}
	}
	fmt.Print("\n\n")

	return &nonce, hash[:], nil
}

// Validates a block's proof of work
func (pow *ProofOfWork) Validate() (*bool, error) {
	var hashInt big.Int

	data, err := pow.prepareData(pow.block.Nonce)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return &isValid, nil
}
