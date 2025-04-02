package core

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"go-burrokuchen/model"
	"go-burrokuchen/utils"
	"io/fs"
	"os"
)

// Wallets stores a collection of wallets
type Wallets struct {
	cfg     *model.Config
	Wallets map[string]*Wallet
}

// NewWallets creates Wallets and retrieves it from a file if it exists
func NewWallets(cfg *model.Config) (*Wallets, error) {
	wallets := Wallets{
		cfg:     cfg,
		Wallets: make(map[string]*Wallet),
	}

	err := wallets.LoadFromFile()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	return &wallets, nil
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() (*string, error) {
	wallet, err := NewWallet(ws.cfg)
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	address, err := wallet.GetAddress()
	if err != nil {
		return nil, utils.CatchErr(err)
	}

	addressStr := string(address)

	ws.Wallets[addressStr] = wallet

	return &addressStr, nil
}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile() error {
	walletFile := ws.cfg.WalletConfig.WalletFile

	if _, err := os.Stat(walletFile); errors.Is(err, fs.ErrNotExist) {
		return nil
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		utils.CatchErr(err)
	}

	var wallets Wallets

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		utils.CatchErr(err)
	}

	ws.Wallets = wallets.Wallets

	return nil
}

// GetWallet returns a Wallet by its address
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// SaveToFile saves wallets to a file
func (ws *Wallets) SaveToFile() error {
	walletFile := ws.cfg.WalletConfig.WalletFile
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		return utils.CatchErr(err)
	}

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		return utils.CatchErr(err)
	}

	return nil
}
