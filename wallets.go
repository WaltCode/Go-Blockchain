package main 

import (
	"bytes"
	"log"
	"os"
	"fmt"
	"encoding/gob"
	"crypto/elliptic"
)

const walletFile = "wallet_%s.dat"

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets(nodeId string) (*Wallets, error) {
	wallets:= Wallets{}
	wallets.Wallets = make(map[string]*Wallet)
	err:= wallets.LoadFromFile(nodeId)
	return &wallets, err
}

func (ws *Wallets) createWallet() string{
	wallet:= NewWallet()
	address:= fmt.Sprintf("%s", wallet.GetAddress())
	ws.Wallets[address] = wallet
	return address
}

func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address:= range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadFromFile(nodeId string) error {
	walletFile:= fmt.Sprintf(walletFile, nodeId)
	if _, err:= os.Stat(walletFile);
	os.IsNotExist(err){
		return err
	}
	fileContent, err:= os.ReadFile(walletFile)
	if err != nil {
		log.Panic((err))
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder:= gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic((err))
	}
	ws.Wallets = wallets.Wallets
	return nil
}

func (ws Wallets) SaveToFile(nodeId string) {
	var content bytes.Buffer
	walletFile:= fmt.Sprintf(walletFile, nodeId)
	gob.Register(elliptic.P256())
	encoder:= gob.NewEncoder(&content)
	err:= encoder.Encode(ws)
	if err != nil {
		log.Panic((err))
	}
	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Panic((err))
	}
}