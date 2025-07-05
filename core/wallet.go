package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

// Wallet stores private and public keys
type Wallet struct {
	cfg        *model.Config
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet(cfg *model.Config) (*Wallet, error) {
	privKey, pubKey, err := newKeyPair()
	if err != nil {
		return nil, utils.CatchErr(err)
	}
	wallet := Wallet{
		cfg:        cfg,
		PrivateKey: *privKey,
		PublicKey:  pubKey,
	}

	return &wallet, nil
}

// newKeyPair generates and returns a private and public key pair
func newKeyPair() (*ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()

	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, utils.CatchErr(err)
	}

	pubKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)

	return privKey, pubKey, nil
}

// HashPubKey hashes public key
func HashPubKey(pubKey []byte) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
}

func checkSum(payload []byte, checkSumLength int) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:checkSumLength]
}

// GetAddress returns wallet address
func (w *Wallet) GetAddress() ([]byte, error) {
	pubKeyHash, err := HashPubKey(w.PublicKey)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	versionPayload := append([]byte{version}, pubKeyHash...)

	checkSumLength := w.cfg.WalletConfig.CheckSumLength

	checkSum := checkSum(versionPayload, checkSumLength)

	fullPayload := append(versionPayload, checkSum...)
	address := utils.Base58Encode(fullPayload)

	return address, nil
}

// ValidateAddress checks if address if valid
func ValidateAddress(cfg *model.Config, address string) (*bool, error) {
	checkSumLength := cfg.WalletConfig.CheckSumLength

	pubKeyHash := utils.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checkSumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checkSumLength]
	targetChecksum := checkSum(append([]byte{version}, pubKeyHash...), checkSumLength)

	result := bytes.Equal(actualChecksum, targetChecksum)

	return &result, nil
}

// GobEncode encodes the public key using the streams of gobs
func (w *Wallet) GobEncode() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode the curve name
	curveName := "P256"
	if err := encoder.Encode(curveName); err != nil {
		return nil, err
	}

	// Encode the private key components
	privateKeyBytes := w.PrivateKey.D.Bytes()
	if err := encoder.Encode(privateKeyBytes); err != nil {
		return nil, err
	}

	// Encode the public key components
	xBytes := w.PrivateKey.PublicKey.X.Bytes()
	yBytes := w.PrivateKey.PublicKey.Y.Bytes()
	if err := encoder.Encode(xBytes); err != nil {
		return nil, err
	}
	if err := encoder.Encode(yBytes); err != nil {
		return nil, err
	}

	// Encode the public key
	if err := encoder.Encode(w.PublicKey); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// GobDecode decodes the public key using the streams of gobs
func (w *Wallet) GobDecode(data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	// Decode the curve name
	var curveName string
	if err := decoder.Decode(&curveName); err != nil {
		return err
	}

	// Set the curve based on the name
	var curve elliptic.Curve
	switch curveName {
	case "P256":
		curve = elliptic.P256()
	default:
		return fmt.Errorf("unsupported curve: %s", curveName)
	}

	// Decode the private key components
	var privateKeyBytes []byte
	if err := decoder.Decode(&privateKeyBytes); err != nil {
		return err
	}

	w.PrivateKey.PublicKey.Curve = curve
	w.PrivateKey.D = new(big.Int).SetBytes(privateKeyBytes)

	// Decode the public key components
	var xBytes, yBytes []byte
	if err := decoder.Decode(&xBytes); err != nil {
		return err
	}
	if err := decoder.Decode(&yBytes); err != nil {
		return err
	}

	w.PrivateKey.PublicKey.X = new(big.Int).SetBytes(xBytes)
	w.PrivateKey.PublicKey.Y = new(big.Int).SetBytes(yBytes)

	// Decode the public key
	if err := decoder.Decode(&w.PublicKey); err != nil {
		return err
	}

	return nil
}
