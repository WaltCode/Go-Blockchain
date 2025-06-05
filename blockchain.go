package main

import (
	"fmt"
	"os"
	"bytes"
	"crypto/ecdsa"
	"log"
	"errors"
	"encoding/hex"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchain_%s.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

type BlockChain struct {
	tip []byte
	db *bolt.DB
}

func createBlockchain(address string, nodeId string) *BlockChain { 
	dbFile := fmt.Sprintf(dbFile, nodeId)
	if dbExists(dbFile) {
		fmt.Println("DB already exist")
		os.Exit(1)
	}
	var tip []byte
	cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
	genesis:= NewGenesisBlock(cbtx)

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b, err:= tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc:= BlockChain{tip, db}
	return &bc
}

func newBlockchain(nodeId string) *BlockChain {
	dbFile:= fmt.Sprintf(dbFile, nodeId)
	if dbExists(dbFile) == false{
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}
	var tip []byte
	db, err:= bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func (tx *bolt.Tx) error {
		b:= tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	bc:= BlockChain{tip, db}
	return &bc
}